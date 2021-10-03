package adaptor

import (
	"context"
	"errors"
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

type GroupModel struct {
	ID         string
	UserID     string
	Name       string
	CreateTime time.Time
}

func (g GroupModel) TableName() string {
	return "report_group"
}

type ReportRepository struct {
	client *gorm.DB
}

func NewReportRepository(dbClient *gorm.DB) ReportRepository {
	return ReportRepository{client: dbClient}
}

func (r ReportRepository) AddReport(ctx context.Context, report domain.Report) error {
	return r.client.Create(&ReportModel{
		ID:         report.ID(),
		GroupID:    report.GroupID(),
		UserID:     report.UserID(),
		CreateTime: report.Time(),
		Content:    report.Content(),
	}).Error
}

type allTypeListReceiver struct {
	ID              string    `gorm:"id"`
	GroupID         string    `gorm:"column:group_id"`
	UserID          string    `gorm:"column:user_id"`
	Content         string    `gorm:"column:content"`
	CreateTime      time.Time `gorm:"column:create_time"`
	GroupName       string    `gorm:"column:group_name"`
	GroupCreateTime time.Time `gorm:"column:group_create_time"`
}

func (r ReportRepository) FindReportWithAllGroup(ctx context.Context, userID string, pageNo, pageSize int) (query.RespReportAllGroupList, error) {
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
		var r allTypeListReceiver
		if err = rawData.Scan(&r); err != nil {
			return query.RespReportAllGroupList{}, err
		}
		list = append(list, r)
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
