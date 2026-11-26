package tag

const (
	eduBiz = "edu"
	// 教育标签分组常量 - 按词义范围从大到小排序
	groupStage         = "stage"         // 学段分组（最大范围：幼儿园、小学、初中、高中、大学）
	groupAcademicYear  = "academicYear"  // 学年分组（时间维度：2024-2025学年）
	groupSemester      = "semester"      // 学期分组（时间维度：第一学期、第二学期）
	groupSubject       = "subject"       // 学科分组（课程分类：语文、数学、英语）
	groupGrade         = "grade"         // 年级分组（教育层级：一年级、二年级等）
	groupVersion       = "version"       // 版本分组（教材版本：人教版、苏教版）
	groupTeachMaterial = "teachMaterial" // 教材分组（教学资源）
	groupChapter       = "chapter"       // 章节分组（教材章节）
	groupKnowledge     = "knowledge"     // 知识点分组（最小粒度）
	groupClass         = "class"         // 班级分组（教学组织）
	groupCourse        = "course"        // 课程分组（具体课程）
	groupTopic         = "topic"         // 专题分组
	groupTopicNode     = "topicNode"     // 专题目录节点分组
)

// EduConverter 教育行业标签转换器 - 基于基础 Converter，提供教育行业特定的便捷方法
// 嵌入 Converter 以复用基础功能，同时提供教育行业专用的方法
type EduConverter struct {
	*Converter
}

// NewEduConverter 创建教育行业标签转换器
// 返回:
//   - *EduConverter: 教育行业转换器实例
//   - error: 如果创建失败，返回错误
func NewEduConverter() (*EduConverter, error) {
	base, err := NewConverter(eduBiz)
	if err != nil {
		return nil, err
	}
	return &EduConverter{Converter: base}, nil
}

// ==================== 学段相关方法（词义范围最大）====================

// GetStages 获取学段标签列表
// 返回:
//   - []string: 学段 ID 列表，如果不存在则返回 nil
func (c *EduConverter) GetStages() []string {
	return c.GetGroupTagIDs(groupStage)
}

// SetStages 设置学段标签（替换现有学段）
// 参数:
//   - stageIDs: 学段 ID 列表
//
// 返回:
//   - error: 如果设置失败，返回错误
func (c *EduConverter) SetStages(stageIDs ...string) error {
	return c.SetGroup(groupStage, stageIDs)
}

// AddStages 添加学段标签（支持添加多个，自动去重）
// 参数:
//   - stageIDs: 要添加的学段 ID 列表（可变参数）
//
// 返回:
//   - error: 如果添加失败，返回错误
func (c *EduConverter) AddStages(stageIDs ...string) error {
	return c.AddGroupTagID(groupStage, stageIDs...)
}

// RemoveStages 移除学段标签（支持移除多个）
// 参数:
//   - stageIDs: 要移除的学段 ID 列表（可变参数）
func (c *EduConverter) RemoveStages(stageIDs ...string) {
	c.RemoveGroupTagID(groupStage, stageIDs...)
}

// ==================== 学年相关方法 ====================

// GetAcademicYears 获取学年标签列表
// 返回:
//   - []string: 学年 ID 列表，如果不存在则返回 nil
func (c *EduConverter) GetAcademicYears() []string {
	return c.GetGroupTagIDs(groupAcademicYear)
}

// SetAcademicYears 设置学年标签（替换现有学年）
// 参数:
//   - academicYearIDs: 学年 ID 列表
//
// 返回:
//   - error: 如果设置失败，返回错误
func (c *EduConverter) SetAcademicYears(academicYearIDs ...string) error {
	return c.SetGroup(groupAcademicYear, academicYearIDs)
}

// AddAcademicYears 添加学年标签（支持添加多个，自动去重）
// 参数:
//   - academicYearIDs: 要添加的学年 ID 列表（可变参数）
//
// 返回:
//   - error: 如果添加失败，返回错误
func (c *EduConverter) AddAcademicYears(academicYearIDs ...string) error {
	return c.AddGroupTagID(groupAcademicYear, academicYearIDs...)
}

// RemoveAcademicYears 移除学年标签（支持移除多个）
// 参数:
//   - academicYearIDs: 要移除的学年 ID 列表（可变参数）
func (c *EduConverter) RemoveAcademicYears(academicYearIDs ...string) {
	c.RemoveGroupTagID(groupAcademicYear, academicYearIDs...)
}

// ==================== 学期相关方法 ====================

// GetSemesters 获取学期标签列表
// 返回:
//   - []string: 学期 ID 列表，如果不存在则返回 nil
func (c *EduConverter) GetSemesters() []string {
	return c.GetGroupTagIDs(groupSemester)
}

// SetSemesters 设置学期标签（替换现有学期）
// 参数:
//   - semesterIDs: 学期 ID 列表
//
// 返回:
//   - error: 如果设置失败，返回错误
func (c *EduConverter) SetSemesters(semesterIDs ...string) error {
	return c.SetGroup(groupSemester, semesterIDs)
}

// AddSemesters 添加学期标签（支持添加多个，自动去重）
// 参数:
//   - semesterIDs: 要添加的学期 ID 列表（可变参数）
//
// 返回:
//   - error: 如果添加失败，返回错误
func (c *EduConverter) AddSemesters(semesterIDs ...string) error {
	return c.AddGroupTagID(groupSemester, semesterIDs...)
}

// RemoveSemesters 移除学期标签（支持移除多个）
// 参数:
//   - semesterIDs: 要移除的学期 ID 列表（可变参数）
func (c *EduConverter) RemoveSemesters(semesterIDs ...string) {
	c.RemoveGroupTagID(groupSemester, semesterIDs...)
}

// ==================== 学科相关方法 ====================

// GetSubjects 获取学科标签列表
// 返回:
//   - []string: 学科 ID 列表，如果不存在则返回 nil
func (c *EduConverter) GetSubjects() []string {
	return c.GetGroupTagIDs(groupSubject)
}

// SetSubjects 设置学科标签（替换现有学科）
// 参数:
//   - subjectIDs: 学科 ID 列表
//
// 返回:
//   - error: 如果设置失败，返回错误
func (c *EduConverter) SetSubjects(subjectIDs ...string) error {
	return c.SetGroup(groupSubject, subjectIDs)
}

// AddSubjects 添加学科标签（支持添加多个，自动去重）
// 参数:
//   - subjectIDs: 要添加的学科 ID 列表（可变参数）
//
// 返回:
//   - error: 如果添加失败，返回错误
func (c *EduConverter) AddSubjects(subjectIDs ...string) error {
	return c.AddGroupTagID(groupSubject, subjectIDs...)
}

// RemoveSubjects 移除学科标签（支持移除多个）
// 参数:
//   - subjectIDs: 要移除的学科 ID 列表（可变参数）
func (c *EduConverter) RemoveSubjects(subjectIDs ...string) {
	c.RemoveGroupTagID(groupSubject, subjectIDs...)
}

// ==================== 年级相关方法 ====================

// GetGrades 获取年级标签列表
// 返回:
//   - []string: 年级 ID 列表，如果不存在则返回 nil
func (c *EduConverter) GetGrades() []string {
	return c.GetGroupTagIDs(groupGrade)
}

// SetGrades 设置年级标签（替换现有年级）
// 参数:
//   - gradeIDs: 年级 ID 列表
//
// 返回:
//   - error: 如果设置失败，返回错误
func (c *EduConverter) SetGrades(gradeIDs ...string) error {
	return c.SetGroup(groupGrade, gradeIDs)
}

// AddGrades 添加年级标签（支持添加多个，自动去重）
// 参数:
//   - gradeIDs: 要添加的年级 ID 列表（可变参数）
//
// 返回:
//   - error: 如果添加失败，返回错误
func (c *EduConverter) AddGrades(gradeIDs ...string) error {
	return c.AddGroupTagID(groupGrade, gradeIDs...)
}

// RemoveGrades 移除年级标签（支持移除多个）
// 参数:
//   - gradeIDs: 要移除的年级 ID 列表（可变参数）
func (c *EduConverter) RemoveGrades(gradeIDs ...string) {
	c.RemoveGroupTagID(groupGrade, gradeIDs...)
}

// ==================== 版本相关方法 ====================

// GetVersions 获取版本标签列表
// 返回:
//   - []string: 版本 ID 列表，如果不存在则返回 nil
func (c *EduConverter) GetVersions() []string {
	return c.GetGroupTagIDs(groupVersion)
}

// SetVersions 设置版本标签（替换现有版本）
// 参数:
//   - versionIDs: 版本 ID 列表
//
// 返回:
//   - error: 如果设置失败，返回错误
func (c *EduConverter) SetVersions(versionIDs ...string) error {
	return c.SetGroup(groupVersion, versionIDs)
}

// AddVersions 添加版本标签（支持添加多个，自动去重）
// 参数:
//   - versionIDs: 要添加的版本 ID 列表（可变参数）
//
// 返回:
//   - error: 如果添加失败，返回错误
func (c *EduConverter) AddVersions(versionIDs ...string) error {
	return c.AddGroupTagID(groupVersion, versionIDs...)
}

// RemoveVersions 移除版本标签（支持移除多个）
// 参数:
//   - versionIDs: 要移除的版本 ID 列表（可变参数）
func (c *EduConverter) RemoveVersions(versionIDs ...string) {
	c.RemoveGroupTagID(groupVersion, versionIDs...)
}

// ==================== 教材相关方法 ====================

// GetTeachMaterials 获取教材标签列表
// 返回:
//   - []string: 教材 ID 列表，如果不存在则返回 nil
func (c *EduConverter) GetTeachMaterials() []string {
	return c.GetGroupTagIDs(groupTeachMaterial)
}

// SetTeachMaterials 设置教材标签（替换现有教材）
// 参数:
//   - teachMaterialIDs: 教材 ID 列表
//
// 返回:
//   - error: 如果设置失败，返回错误
func (c *EduConverter) SetTeachMaterials(teachMaterialIDs ...string) error {
	return c.SetGroup(groupTeachMaterial, teachMaterialIDs)
}

// AddTeachMaterials 添加教材标签（支持添加多个，自动去重）
// 参数:
//   - teachMaterialIDs: 要添加的教材 ID 列表（可变参数）
//
// 返回:
//   - error: 如果添加失败，返回错误
func (c *EduConverter) AddTeachMaterials(teachMaterialIDs ...string) error {
	return c.AddGroupTagID(groupTeachMaterial, teachMaterialIDs...)
}

// RemoveTeachMaterials 移除教材标签（支持移除多个）
// 参数:
//   - teachMaterialIDs: 要移除的教材 ID 列表（可变参数）
func (c *EduConverter) RemoveTeachMaterials(teachMaterialIDs ...string) {
	c.RemoveGroupTagID(groupTeachMaterial, teachMaterialIDs...)
}

// ==================== 章节相关方法 ====================

// GetChapters 获取章节标签列表
// 返回:
//   - []string: 章节 ID 列表，如果不存在则返回 nil
func (c *EduConverter) GetChapters() []string {
	return c.GetGroupTagIDs(groupChapter)
}

// SetChapters 设置章节标签（替换现有章节）
// 参数:
//   - chapterIDs: 章节 ID 列表
//
// 返回:
//   - error: 如果设置失败，返回错误
func (c *EduConverter) SetChapters(chapterIDs ...string) error {
	return c.SetGroup(groupChapter, chapterIDs)
}

// AddChapters 添加章节标签（支持添加多个，自动去重）
// 参数:
//   - chapterIDs: 要添加的章节 ID 列表（可变参数）
//
// 返回:
//   - error: 如果添加失败，返回错误
func (c *EduConverter) AddChapters(chapterIDs ...string) error {
	return c.AddGroupTagID(groupChapter, chapterIDs...)
}

// RemoveChapters 移除章节标签（支持移除多个）
// 参数:
//   - chapterIDs: 要移除的章节 ID 列表（可变参数）
func (c *EduConverter) RemoveChapters(chapterIDs ...string) {
	c.RemoveGroupTagID(groupChapter, chapterIDs...)
}

// ==================== 知识点相关方法 ====================

// GetKnowledge 获取知识点标签列表
// 返回:
//   - []string: 知识点 ID 列表，如果不存在则返回 nil
func (c *EduConverter) GetKnowledge() []string {
	return c.GetGroupTagIDs(groupKnowledge)
}

// SetKnowledge 设置知识点标签（替换现有知识点）
// 参数:
//   - knowledgeIDs: 知识点 ID 列表
//
// 返回:
//   - error: 如果设置失败，返回错误
func (c *EduConverter) SetKnowledge(knowledgeIDs ...string) error {
	return c.SetGroup(groupKnowledge, knowledgeIDs)
}

// AddKnowledge 添加知识点标签（支持添加多个，自动去重）
// 参数:
//   - knowledgeIDs: 要添加的知识点 ID 列表（可变参数）
//
// 返回:
//   - error: 如果添加失败，返回错误
func (c *EduConverter) AddKnowledge(knowledgeIDs ...string) error {
	return c.AddGroupTagID(groupKnowledge, knowledgeIDs...)
}

// RemoveKnowledge 移除知识点标签（支持移除多个）
// 参数:
//   - knowledgeIDs: 要移除的知识点 ID 列表（可变参数）
func (c *EduConverter) RemoveKnowledge(knowledgeIDs ...string) {
	c.RemoveGroupTagID(groupKnowledge, knowledgeIDs...)
}

// ==================== 班级相关方法 ====================

// GetClasses 获取班级标签列表
// 返回:
//   - []string: 班级 ID 列表，如果不存在则返回 nil
func (c *EduConverter) GetClasses() []string {
	return c.GetGroupTagIDs(groupClass)
}

// SetClasses 设置班级标签（替换现有班级）
// 参数:
//   - classIDs: 班级 ID 列表
//
// 返回:
//   - error: 如果设置失败，返回错误
func (c *EduConverter) SetClasses(classIDs ...string) error {
	return c.SetGroup(groupClass, classIDs)
}

// AddClasses 添加班级标签（支持添加多个，自动去重）
// 参数:
//   - classIDs: 要添加的班级 ID 列表（可变参数）
//
// 返回:
//   - error: 如果添加失败，返回错误
func (c *EduConverter) AddClasses(classIDs ...string) error {
	return c.AddGroupTagID(groupClass, classIDs...)
}

// RemoveClasses 移除班级标签（支持移除多个）
// 参数:
//   - classIDs: 要移除的班级 ID 列表（可变参数）
func (c *EduConverter) RemoveClasses(classIDs ...string) {
	c.RemoveGroupTagID(groupClass, classIDs...)
}

// ==================== 课程相关方法（词义范围最小）====================

// GetCourses 获取课程标签列表
// 返回:
//   - []string: 课程 ID 列表，如果不存在则返回 nil
func (c *EduConverter) GetCourses() []string {
	return c.GetGroupTagIDs(groupCourse)
}

// SetCourses 设置课程标签（替换现有课程）
// 参数:
//   - courseIDs: 课程 ID 列表
//
// 返回:
//   - error: 如果设置失败，返回错误
func (c *EduConverter) SetCourses(courseIDs ...string) error {
	return c.SetGroup(groupCourse, courseIDs)
}

// AddCourses 添加课程标签（支持添加多个，自动去重）
// 参数:
//   - courseIDs: 要添加的课程 ID 列表（可变参数）
//
// 返回:
//   - error: 如果添加失败，返回错误
func (c *EduConverter) AddCourses(courseIDs ...string) error {
	return c.AddGroupTagID(groupCourse, courseIDs...)
}

// RemoveCourses 移除课程标签（支持移除多个）
// 参数:
//   - courseIDs: 要移除的课程 ID 列表（可变参数）
func (c *EduConverter) RemoveCourses(courseIDs ...string) {
	c.RemoveGroupTagID(groupCourse, courseIDs...)
}

// ==================== 专题相关方法 ====================

// GetTopics 获取专题 ID 列表
func (c *EduConverter) GetTopics() []string {
	return c.GetGroupTagIDs(groupTopic)
}

// SetTopics 设置专题 ID 列表（替换现有专题）
func (c *EduConverter) SetTopics(topicIDs ...string) error {
	return c.SetGroup(groupTopic, topicIDs)
}

// AddTopics 添加专题 ID
func (c *EduConverter) AddTopics(topicIDs ...string) error {
	return c.AddGroupTagID(groupTopic, topicIDs...)
}

// RemoveTopics 移除专题 ID
func (c *EduConverter) RemoveTopics(topicIDs ...string) {
	c.RemoveGroupTagID(groupTopic, topicIDs...)
}

// GetTopicNodes 获取专题目录节点 ID 列表
func (c *EduConverter) GetTopicNodes() []string {
	return c.GetGroupTagIDs(groupTopicNode)
}

// SetTopicNodes 设置专题目录节点 ID 列表（替换现有节点）
func (c *EduConverter) SetTopicNodes(nodeIDs ...string) error {
	return c.SetGroup(groupTopicNode, nodeIDs)
}

// AddTopicNodes 添加专题目录节点 ID
func (c *EduConverter) AddTopicNodes(nodeIDs ...string) error {
	return c.AddGroupTagID(groupTopicNode, nodeIDs...)
}

// RemoveTopicNodes 移除专题目录节点 ID
func (c *EduConverter) RemoveTopicNodes(nodeIDs ...string) {
	c.RemoveGroupTagID(groupTopicNode, nodeIDs...)
}
