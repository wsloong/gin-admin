package schema

import (
	"strings"
	"time"
)

// Menu 菜单对象
type Menu struct {
	RecordID   string      `json:"record_id"`                                  // 记录ID
	Name       string      `json:"name" binding:"required"`                    // 菜单名称
	Sequence   int         `json:"sequence"`                                   // 排序值
	Icon       string      `json:"icon"`                                       // 菜单图标
	Router     string      `json:"router"`                                     // 访问路由
	ParentID   string      `json:"parent_id"`                                  // 父级ID
	ParentPath string      `json:"parent_path"`                                // 父级路径
	ShowStatus int         `json:"show_status" binding:"required,max=2,min=1"` // 显示状态(1:显示 2:隐藏)
	Status     int         `json:"status" binding:"required,max=2,min=1"`      // 状态(1:启用 2:禁用)
	Memo       string      `json:"memo"`                                       // 备注
	Creator    string      `json:"creator"`                                    // 创建者
	CreatedAt  time.Time   `json:"created_at"`                                 // 创建时间
	UpdatedAt  time.Time   `json:"updated_at"`                                 // 更新时间
	Actions    MenuActions `json:"actions"`                                    // 动作列表
}

// MenuQueryParam 查询条件
type MenuQueryParam struct {
	RecordIDs        []string `form:"-"`          // 记录ID列表
	Name             string   `form:"-"`          // 菜单名称
	PrefixParentPath string   `form:"-"`          // 父级路径(前缀模糊查询)
	LikeName         string   `form:"likeName"`   // 菜单名称(模糊查询)
	ParentID         *string  `form:"parentID"`   // 父级内码
	ShowStatus       int      `json:"showStatus"` // 显示状态(1:显示 2:隐藏)
	Status           int      `form:"status"`     // 状态(1:启用 2:禁用)
}

// MenuQueryOptions 查询可选参数项
type MenuQueryOptions struct {
	PageParam   *PaginationParam // 分页参数
	OrderFields []*OrderField    // 排序字段
}

// MenuQueryResult 查询结果
type MenuQueryResult struct {
	Data       Menus
	PageResult *PaginationResult
}

// Menus 菜单列表
type Menus []*Menu

// ToMap 转换为键值映射
func (a Menus) ToMap() map[string]*Menu {
	m := make(map[string]*Menu)
	for _, item := range a {
		m[item.RecordID] = item
	}
	return m
}

// SplitAndGetAllRecordIDs 拆分父级路径并获取所有记录ID
func (a Menus) SplitAndGetAllRecordIDs() []string {
	recordIDs := make([]string, 0, len(a))
	for _, item := range a {
		recordIDs = append(recordIDs, item.RecordID)
		if item.ParentPath == "" {
			continue
		}

		pps := strings.Split(item.ParentPath, "/")
		for _, pp := range pps {
			var exists bool
			for _, recordID := range recordIDs {
				if pp == recordID {
					exists = true
					break
				}
			}
			if !exists {
				recordIDs = append(recordIDs, pp)
			}
		}
	}
	return recordIDs
}

// ToTree 转换为菜单树
func (a Menus) ToTree() MenuTrees {
	list := make(MenuTrees, len(a))
	for i, item := range a {
		list[i] = &MenuTree{
			RecordID:   item.RecordID,
			Name:       item.Name,
			Icon:       item.Icon,
			Router:     item.Router,
			ParentID:   item.ParentID,
			ParentPath: item.ParentPath,
			ShowStatus: item.ShowStatus,
			Actions:    item.Actions,
		}
	}
	return list.ToTree()
}

// ToLeafRecordIDs 转换为叶子节点记录ID列表
func (a Menus) ToLeafRecordIDs() []string {
	var leafNodeIDs []string
	tree := a.ToTree()
	a.fillLeafNodeID(&tree, &leafNodeIDs)
	return leafNodeIDs
}

func (a Menus) fillLeafNodeID(tree *MenuTrees, leafNodeIDs *[]string) {
	for _, node := range *tree {
		if node.Children == nil || len(*node.Children) == 0 {
			*leafNodeIDs = append(*leafNodeIDs, node.RecordID)
			continue
		}
		a.fillLeafNodeID(node.Children, leafNodeIDs)
	}
}

// ----------------------------------------MenuTree--------------------------------------

// MenuTree 菜单树
type MenuTree struct {
	RecordID   string      `json:"record_id"`          // 记录ID
	Name       string      `json:"name"`               // 菜单名称
	Icon       string      `json:"icon"`               // 菜单图标
	Router     string      `json:"router"`             // 访问路由
	ParentID   string      `json:"parent_id"`          // 父级ID
	ParentPath string      `json:"parent_path"`        // 父级路径
	ShowStatus int         `json:"show_status"`        // 显示状态(1:显示 2:隐藏)
	Actions    MenuActions `json:"actions"`            // 动作列表
	Children   *MenuTrees  `json:"children,omitempty"` // 子级树
}

// MenuTrees 菜单树列表
type MenuTrees []*MenuTree

// ToTree 转换为树形结构
func (a MenuTrees) ToTree() []*MenuTree {
	mi := make(map[string]*MenuTree)
	for _, item := range a {
		mi[item.RecordID] = item
	}

	var list []*MenuTree
	for _, item := range a {
		if item.ParentID == "" {
			list = append(list, item)
			continue
		}
		if pitem, ok := mi[item.ParentID]; ok {
			if pitem.Children == nil {
				var children MenuTrees
				children = append(children, item)
				pitem.Children = &children
				continue
			}
			*pitem.Children = append(*pitem.Children, item)
		}
	}
	return list
}

// ----------------------------------------MenuAction--------------------------------------

// MenuAction 菜单动作对象
type MenuAction struct {
	RecordID  string              `json:"record_id"`                  // 记录ID
	MenuID    string              `json:"menu_id" binding:"required"` // 菜单ID
	Code      string              `json:"code" binding:"required"`    // 动作编号
	Name      string              `json:"name" binding:"required"`    // 动作名称
	Resources MenuActionResources `json:"resources"`                  // 资源列表
}

// MenuActionQueryParam 查询条件
type MenuActionQueryParam struct {
	MenuID string // 菜单ID
}

// MenuActionQueryOptions 查询可选参数项
type MenuActionQueryOptions struct {
	PageParam   *PaginationParam // 分页参数
	OrderFields []*OrderField    // 排序字段
}

// MenuActionQueryResult 查询结果
type MenuActionQueryResult struct {
	Data       MenuActions
	PageResult *PaginationResult
}

// MenuActions 菜单动作管理列表
type MenuActions []*MenuAction

// FillResources 填充资源数据
func (a MenuActions) FillResources(mResources map[string]MenuActionResources) {
	for i, item := range a {
		a[i].Resources = mResources[item.RecordID]
	}
}

// GetByRecordID 根据记录ID获取数据项
func (a MenuActions) GetByRecordID(recordID string) *MenuAction {
	for _, item := range a {
		if item.RecordID == recordID {
			return item
		}
	}
	return nil
}

// ----------------------------------------MenuActionResource--------------------------------------

// MenuActionResource 菜单动作关联资源对象
type MenuActionResource struct {
	RecordID string `json:"record_id"`                    // 记录ID
	ActionID string `json:"action_id" binding:"required"` // 菜单动作ID
	Method   string `json:"method" binding:"required"`    // 资源请求方式(支持正则)
	Path     string `json:"path" binding:"required"`      // 资源请求路径（支持/:id匹配）
}

// MenuActionResourceQueryParam 查询条件
type MenuActionResourceQueryParam struct {
	MenuID string // 菜单ID
}

// MenuActionResourceQueryOptions 查询可选参数项
type MenuActionResourceQueryOptions struct {
	PageParam   *PaginationParam // 分页参数
	OrderFields []*OrderField    // 排序字段
}

// MenuActionResourceQueryResult 查询结果
type MenuActionResourceQueryResult struct {
	Data       MenuActionResources
	PageResult *PaginationResult
}

// MenuActionResources 菜单动作关联资源管理列表
type MenuActionResources []*MenuActionResource

// ToActionIDMap 转换为动作ID映射
func (a MenuActionResources) ToActionIDMap() map[string]MenuActionResources {
	m := make(map[string]MenuActionResources)
	for _, item := range a {
		if v, ok := m[item.ActionID]; ok {
			v = append(v, item)
			m[item.ActionID] = v
			continue
		}
		m[item.ActionID] = MenuActionResources{item}
	}
	return m
}

// GetByRecordID 根据记录ID获取数据项
func (a MenuActionResources) GetByRecordID(recordID string) *MenuActionResource {
	for _, item := range a {
		if item.RecordID == recordID {
			return item
		}
	}
	return nil
}
