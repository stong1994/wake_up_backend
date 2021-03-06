package adaptor

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"sort"
	"wake_up_backend/internal/rethink/app/query"
	"wake_up_backend/internal/rethink/domain"
)

type RethinkRepo struct {
	client *mongo.Database
}

func NewRethinkRepo(client *mongo.Database) *RethinkRepo {
	return &RethinkRepo{
		client: client,
	}
}

type RethinkModel struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	UserID        string             `bson:"user_id"`
	ReportTime    primitive.DateTime `bson:"report_time"`
	ReportGroupID primitive.ObjectID `bson:"report_group_id"`

	Content    string             `bson:"content"`
	RecordTime primitive.DateTime `bson:"record_time"`

	RethinkContent string             `bson:"rethink_content"`
	RethinkTime    primitive.DateTime `bson:"rethink_time"`
}

func (RethinkModel) CollectionName() string {
	return "rethink"
}

func (r RethinkRepo) AddReport(ctx context.Context, report domain.Report) (string, error) {
	groupID, err := getObjID(report.GroupID())
	if err != nil {
		return "", err
	}
	resp, err := r.client.Collection(RethinkModel{}.CollectionName()).InsertOne(ctx, &RethinkModel{
		ReportGroupID: groupID,
		UserID:        report.UserID(),
		ReportTime:    primitive.NewDateTimeFromTime(report.Time()),
		Content:       report.Content(),
	})
	if err != nil {
		return "", err
	}
	return resp.InsertedID.(primitive.ObjectID).String(), nil
}

func (r RethinkRepo) FindRethinkByID(ctx context.Context, rethinkID string) (domain.Rethink, error) {
	id, err := getObjID(rethinkID)
	if err != nil {
		return domain.Rethink{}, err
	}
	cond := bson.M{
		"_id": id,
	}
	var model RethinkModel
	if err = r.client.Collection(RethinkModel{}.CollectionName()).FindOne(ctx, cond).Decode(&model); err != nil {
		return domain.Rethink{}, err
	}
	return domain.UnmarshalRethinkFromDB(model.ID.Hex(), model.UserID, model.ReportGroupID.Hex(), model.Content,
		model.RethinkContent, model.ReportTime.Time(), model.RecordTime.Time(), model.RethinkTime.Time()), nil
}

func (r RethinkRepo) SaveRethink(ctx context.Context, rethink domain.Rethink) error {
	id, err := getObjID(rethink.ID())
	if err != nil {
		return err
	}
	groupID, err := getObjID(rethink.GroupID())
	if err != nil {
		return err
	}
	cond := bson.M{
		"_id": id,
	}
	_, err = r.client.Collection(RethinkModel{}.CollectionName()).ReplaceOne(ctx, cond, &RethinkModel{
		//ID:             id,
		UserID:         rethink.UserID(),
		ReportTime:     primitive.NewDateTimeFromTime(rethink.ReportTime()),
		ReportGroupID:  groupID,
		Content:        rethink.Content(),
		RecordTime:     primitive.NewDateTimeFromTime(rethink.RecordTime()),
		RethinkContent: rethink.RethinkContent(),
		RethinkTime:    primitive.NewDateTimeFromTime(rethink.RethinkTime()),
	})
	return err
}

type ReportGroupModelV2 struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name"`
	UserID     string             `bson:"user_id"`
	CreateTime primitive.DateTime `bson:"create_time"`
}

func (ReportGroupModelV2) CollectionName() string {
	return "rethink_group"
}

func (r RethinkRepo) AddReportGroup(ctx context.Context, reportGroup domain.ReportGroup) (string, error) {
	resp, err := r.client.Collection(ReportGroupModelV2{}.CollectionName()).InsertOne(ctx, &ReportGroupModelV2{
		UserID:     reportGroup.UserID(),
		CreateTime: primitive.NewDateTimeFromTime(reportGroup.CreateTime()),
		Name:       reportGroup.Name(),
	})
	if err != nil {
		return "", err
	}
	return resp.InsertedID.(primitive.ObjectID).String(), nil
}

func (r RethinkRepo) FindReportGroups(ctx context.Context, userID string) (query.RespReportGroupList, error) {
	matchStage := bson.M{
		"user_id": userID,
	}

	cursor, err := r.client.Collection(ReportGroupModelV2{}.CollectionName()).Find(ctx, matchStage)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var groups []ReportGroupModelV2
	if err = cursor.All(ctx, &groups); err != nil {
		return nil, err
	}

	groupRethinkMap, err := r.getUserGroupRethinkMap(ctx, userID)
	if err != nil {
		return nil, err
	}

	return convert2ReportGroups(groups, groupRethinkMap), nil
}

func convert2ReportGroups(groups []ReportGroupModelV2, rethinkMap map[string]int) query.RespReportGroupList {
	result := make([]query.RespReportGroupItem, len(groups))
	for i, v := range groups {
		result[i] = query.RespReportGroupItem{
			GroupID: v.ID.Hex(),
			Name:    v.Name,
			Count:   rethinkMap[v.ID.Hex()],
		}
	}
	return result
}

type userGroupRethinkMapInfo struct {
	ID    primitive.ObjectID `bson:"_id"`
	Count int                `bson:"count"`
}

func (r RethinkRepo) getUserGroupRethinkMap(ctx context.Context, userID string) (map[string]int, error) {
	matchStage := bson.D{
		{"$match", bson.M{
			"user_id": userID,
		}},
	}
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$group_id"},
			{"count", bson.M{"$sum": 1}},
		}},
	}
	pip := mongo.Pipeline{matchStage, groupStage}
	cursor, err := r.client.Collection(RethinkModel{}.CollectionName()).Aggregate(ctx, pip)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rawMap []userGroupRethinkMapInfo
	if err = cursor.All(ctx, &rawMap); err != nil {
		return nil, err
	}

	result := make(map[string]int)
	for _, v := range rawMap {
		result[v.ID.Hex()] = v.Count
	}
	return result, nil
}

type reportWithGroupInfo struct {
	RethinkModel `bson:"inline"`    // mongo???decoder????????????inline tag??????????????????field????????????
	Group        ReportGroupModelV2 `bson:"group_info"`
}

// todo sort
func (r RethinkRepo) FindUserReports(ctx context.Context, userID string, pageNo, pageSize int) (
	result query.RespReportAllGroupList, err error) {

	totalReportCount, err := r.getUserReportCount(ctx, userID)
	if err != nil {
		return result, err
	}
	matchStage := bson.D{
		{"$match", bson.M{
			"user_id": userID,
		}},
	}
	lookupStage := bson.D{
		{"$lookup", bson.M{
			"from":         ReportGroupModelV2{}.CollectionName(),
			"localField":   "report_group_id",
			"foreignField": "_id",
			"as":           "group_info",
		}},
	}
	unwindStage := bson.D{
		{
			"$unwind", bson.M{
				"path":                       "$group_info",
				"preserveNullAndEmptyArrays": true,
			},
		},
	}
	limitStage := bson.D{
		{
			"$limit", pageSize,
		},
	}
	skipStage := bson.D{
		{
			"$skip", (pageNo - 1) * pageSize,
		},
	}

	pip := mongo.Pipeline{matchStage, lookupStage, unwindStage, limitStage, skipStage}
	cursor, err := r.client.Collection(RethinkModel{}.CollectionName()).Aggregate(ctx, pip)
	if err != nil {
		return result, err
	}
	defer cursor.Close(ctx)

	var rawData []reportWithGroupInfo
	if err = cursor.All(ctx, &rawData); err != nil {
		return result, err
	}

	return query.RespReportAllGroupList{
		Total: totalReportCount,
		List:  convert2Reports(rawData),
	}, nil
}

func convert2Reports(data []reportWithGroupInfo) []query.RespReportAllGroupItem {
	result := make([]query.RespReportAllGroupItem, 0)

	reportsMap := make(map[string][]query.RespReportSingleTypeItem)
	groupMap := make(map[string]string)

	for _, v := range data {
		if _, exist := reportsMap[v.Group.ID.Hex()]; exist {
			reportsMap[v.Group.ID.Hex()] = append(reportsMap[v.Group.ID.Hex()], query.RespReportSingleTypeItem{
				RethinkID:           v.ID.Hex(),
				ReportContent:       v.Content,
				ReportTime:          v.ReportTime.Time(),
				RethinkShortContent: v.RethinkContent,
			})
			continue
		}
		reportsMap[v.Group.ID.Hex()] = []query.RespReportSingleTypeItem{{
			RethinkID:           v.ID.Hex(),
			ReportContent:       v.Content,
			ReportTime:          v.ReportTime.Time(),
			RethinkShortContent: v.RethinkContent,
		}}
		groupMap[v.Group.ID.Hex()] = v.Group.Name
	}
	for groupID, groupName := range groupMap {
		reports := reportsMap[groupID]
		sort.Slice(reports, func(i, j int) bool {
			return reports[i].ReportTime.After(reports[j].ReportTime)
		})
		result = append(result, query.RespReportAllGroupItem{
			GroupID:   groupID,
			GroupName: groupName,
			List:      reports,
		})
	}
	// ??????????????????????????????
	//sort.Slice(result, func(i, j int) bool {
	//	return result[i].
	//})
	return result
}

func (r RethinkRepo) getUserReportCount(ctx context.Context, userID string) (result int64, err error) {
	cond := bson.M{"user_id": userID}
	return r.client.Collection(RethinkModel{}.CollectionName()).CountDocuments(ctx, cond)
}

func (r RethinkRepo) CheckGroup(ctx context.Context, userID, groupID string) (bool, error) {
	gID, err := getObjID(groupID)
	if err != nil {
		return false, err
	}
	matchStage := bson.D{
		{"_id", gID},
		{"user_id", userID},
	}
	bytes, err := r.client.Collection(ReportGroupModelV2{}.CollectionName()).FindOne(ctx, matchStage).DecodeBytes()
	if err != nil {
		return false, err
	}
	vs, err := bytes.Values()
	return len(vs) > 0, err
}

func (r RethinkRepo) FindAllReport(ctx context.Context, userID string, pageNo, pageSize int) ([]query.AllReport, error) {
	matchStage := bson.D{
		{"$match", bson.D{
			{"user_id", userID},
			{"content", bson.M{"$ne": ""}},
		}},
	}
	lookupStage := bson.D{
		{"$lookup", bson.M{
			"from":         ReportGroupModelV2{}.CollectionName(),
			"localField":   "report_group_id",
			"foreignField": "_id",
			"as":           "group_info",
		}},
	}
	unwindStage := bson.D{
		{
			"$unwind", bson.M{
				"path":                       "$group_info",
				"preserveNullAndEmptyArrays": true,
			},
		},
	}
	limitStage := bson.D{
		{
			"$limit", pageSize,
		},
	}
	skipStage := bson.D{
		{
			"$skip", (pageNo - 1) * pageSize,
		},
	}

	pip := mongo.Pipeline{matchStage, lookupStage, unwindStage, limitStage, skipStage}
	cursor, err := r.client.Collection(RethinkModel{}.CollectionName()).Aggregate(ctx, pip)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rawData []reportWithGroupInfo
	if err = cursor.All(ctx, &rawData); err != nil {
		return nil, err
	}

	return convert2AllReport(rawData), nil
}

func convert2AllReport(data []reportWithGroupInfo) []query.AllReport {
	result := make([]query.AllReport, len(data))
	for i, v := range data {
		result[i] = query.AllReport{
			ID:         v.ID.Hex(),
			Content:    v.Content,
			ReportTime: v.ReportTime.Time(),
			GroupID:    v.Group.ID.Hex(),
			GroupName:  v.Group.Name,
		}
	}
	return result
}

func getObjID(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}
