package adaptor

import (
	"context"
	"errors"
	"github.com/stong1994/kit_golang/sstr"
	"gorm.io/gorm"
	"time"
	"wake_up_backend/internal/rethink/app/query"
	"wake_up_backend/internal/rethink/domain"
)

type ReportModel struct {
	ID         string
	GroupID    string
	UserID     string
	CreateTime time.Time
	Content    string
}

func (r ReportModel) TableName() string {
	return "report"
}

type ReportGroupModel struct {
	ID         string
	UserID     string
	Name       string
	CreateTime time.Time
}

func (g ReportGroupModel) TableName() string {
	return "report_group"
}

type RethinkModule struct {
	ID         string
	ReportID   string
	Content    string
	CreateTime time.Time
}

func (r RethinkModule) TableName() string {
	return "rethink"
}

type ReportRepository struct {
	client *gorm.DB
}

func NewReportRepository(dbClient *gorm.DB) ReportRepository {
	return ReportRepository{client: dbClient}
}

func (r ReportRepository) AddReport(ctx context.Context, report domain.Report) (string, error) {
	id := sstr.UUIDHex()
	return id, r.client.Create(&ReportModel{
		ID:         id,
		GroupID:    report.GroupID(),
		UserID:     report.UserID(),
		CreateTime: report.Time(),
		Content:    report.Content(),
	}).Error
}

func (r ReportRepository) AddReportGroup(ctx context.Context, reportGroup domain.ReportGroup) (string, error) {
	id := sstr.UUIDHex()
	return id, r.client.Create(&ReportGroupModel{
		ID:         id,
		UserID:     reportGroup.UserID(),
		CreateTime: reportGroup.CreateTime(),
		Name:       reportGroup.Name(),
	}).Error
}

type allTypeListReceiver struct {
	ID              string    `gorm:"column:id"`
	GroupID         string    `gorm:"column:group_id"`
	UserID          string    `gorm:"column:user_id"`
	Content         string    `gorm:"column:content"`
	CreateTime      time.Time `gorm:"column:create_time"`
	GroupName       string    `gorm:"column:group_name"`
	GroupCreateTime time.Time `gorm:"column:group_create_time"`
}

func (r ReportRepository) FindUserReports(ctx context.Context, userID string, pageNo, pageSize int) (query.RespReportAllGroupList, error) {
	offset := (pageNo - 1) * pageSize
	if pageSize < 0 {
		return query.RespReportAllGroupList{}, errors.New("page_size must be over than 0")
	}
	if offset < 0 {
		return query.RespReportAllGroupList{}, errors.New("page_no must be over than 0")
	}
	count, err := r.getUserAllReportCount(ctx, userID)
	if err != nil {
		return query.RespReportAllGroupList{}, err
	}

	rawSql := `
SELECT
	A.id,
	A.group_id,
	A.user_id,
	A.content,
	A.create_time,
	B.name AS group_name,
	B.create_time AS group_create_time 
FROM
	report_group B
	INNER JOIN report A ON A.group_id = B.id 
	AND A.user_id = ? 
ORDER BY
	B.create_time DESC,
	A.create_time DESC 
	LIMIT ? OFFSET ?;`
	rawData, err := r.client.Raw(rawSql, userID, pageSize, offset).Rows()
	if err != nil {
		return query.RespReportAllGroupList{}, err
	}
	defer func() {
		// panic handle
		if err != nil {
			if e := rawData.Close(); e != nil {
				// log
			}
		}
	}()
	var list []allTypeListReceiver
	for rawData.Next() {
		var rec allTypeListReceiver
		if err = r.client.ScanRows(rawData, &rec); err != nil {
			return query.RespReportAllGroupList{}, err
		}
		list = append(list, rec)
	}
	return r.allTypeListToQuery(count, list), nil
}

func (r ReportRepository) allTypeListToQuery(count int64, data []allTypeListReceiver) query.RespReportAllGroupList {
	var (
		list        []query.RespReportAllGroupItem
		lastGroupID string
	)
	var subList []query.RespReportSingleTypeItem
	for _, v := range data {
		item := query.RespReportSingleTypeItem{
			Content:    v.Content,
			ReportTime: v.CreateTime,
			//RethinkShortContent: "",
			//RethinkContentID:    "",
		}

		if lastGroupID == v.GroupID {
			subList = append(subList, item)
			continue
		}
		list = append(list, query.RespReportAllGroupItem{
			GroupID:   v.GroupID,
			GroupName: v.GroupName,
			List:      subList,
		})
		subList = []query.RespReportSingleTypeItem{item}
	}
	if len(subList) > 0 {
		list = append(list, query.RespReportAllGroupItem{
			GroupID:   data[len(data)-1].GroupID,
			GroupName: data[len(data)-1].GroupName,
			List:      subList,
		})
	}
	return query.RespReportAllGroupList{
		Total: count,
		List:  list,
	}
}

func (r ReportRepository) getUserAllReportCount(ctx context.Context, userID string) (count int64, err error) {
	err = r.client.Model(ReportModel{}).Where("user_id = ?", userID).Count(&count).Error
	return
}

func (r ReportRepository) CheckGroup(ctx context.Context, userID, groupID string) (pass bool, err error) {
	var count int64
	if err = r.client.Model(ReportGroupModel{}).Where("id = ? AND user_id = ?", groupID, userID).
		Count(&count).Error; err != nil {
		return
	}
	if count >= 1 {
		pass = true
	}
	return
}

func (r ReportRepository) FindReportGroups(ctx context.Context, userID string) (query.RespReportGroupList, error) {
	var (
		id     string
		name   string
		cnt    int
		err    error
		result query.RespReportGroupList
	)

	rawSql := `
SELECT A.id, A.name, IFNULL(B.cnt,0) as cnt
FROM report_group A
LEFT JOIN (
	SELECT group_id, COUNT(1) AS cnt 
	FROM report 
	WHERE user_id = ?
	GROUP BY group_id
) B ON A.id = B.group_id
WHERE A.user_id = ?
ORDER BY create_time DESC;
` // TODO opt WHERE
	rows, err := r.client.Raw(rawSql, userID, userID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&id, &name, &cnt); err != nil {
			return nil, err
		}
		result = append(result, query.RespReportGroupItem{
			GroupID: id,
			Name:    name,
			Count:   cnt,
		})
	}
	return result, nil
}

func (r ReportRepository) FindAllReport(ctx context.Context, userID string, pageNo, pageSize int) ([]query.AllReport, error) {
	var (
		id         string
		content    string
		createTime time.Time
		groupID    string
		groupName  string
		err        error
		result     []query.AllReport
	)

	rawSql := `
SELECT A.id, A.content, A.create_time, A.group_id, B.name AS group_name
FROM report A LEFT JOIN report_group B ON A.group_id = B.id
WHERE A.user_id = ?
ORDER BY A.create_time DESC
LIMIT ? OFFSET ?;
` // TODO opt WHERE
	rows, err := r.client.Raw(rawSql, userID, pageSize, (pageNo-1)*pageSize).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&id, &content, &createTime, &groupID, &groupName); err != nil {
			return nil, err
		}
		result = append(result, query.AllReport{
			ID:         id,
			GroupID:    groupID,
			Content:    content,
			ReportTime: createTime,
			GroupName:  groupName,
		})
	}
	return result, nil
}
