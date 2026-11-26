// Package dict service 用于标识不同的微服务
/* cSpell:disable */
package dict

import "github.com/zhixingdx/oipivot/backend/pkg/lang"

// Service 服务
type Service string

const (
	// ServiceBaseTenant 平台租户服务
	ServiceBaseTenant Service = "base-tenant"

	// ServiceBaseIam 基础-身份认证服务
	ServiceBaseIam Service = "base-iam"
	// ServiceBaseNotify 基础-通知服务
	ServiceBaseNotify Service = "base-notify"
	// ServiceBaseSafe 基础-安全服务
	ServiceBaseSafe Service = "base-safe"
	// ServiceBaseWechat 基础-微信服务
	ServiceBaseWechat Service = "base-wechat"
	// ServiceBaseStatic 基础-静态服务
	ServiceBaseStatic Service = "base-static"
	// ServiceBaseStat 基础-统计服务
	ServiceBaseStat Service = "base-stat"
	// ServiceBaseUser 基础-用户服务
	ServiceBaseUser Service = "base-user"
	// ServiceBaseRTC
	ServiceBaseRTC Service = "base-rtc"
	// ServiceBaseCRDT 基础-CRDT服务
	ServiceBaseCRDT Service = "base-crdt"
	// ServiceBaseWorkflow 基础-工作流服务
	ServiceBaseWorkflow Service = "base-workflow"
	// ServiceBaseAudit 基础-审计服务
	ServiceBaseAudit Service = "base-audit"
	// ServiceBaseRag 基础-RAG 平台服务
	ServiceBaseRag Service = "base-rag"
	// ServiceBaseAnalysis 基础-通用分析平台服务
	ServiceBaseAnalysis Service = "base-analysis"
	// ServiceBaseSubscription 基础-订阅与权益服务
	ServiceBaseSubscription Service = "base-subscription"
	// ServiceBaseArtifact 基础-文件产物处理平台服务
	ServiceBaseArtifact Service = "base-artifact"
	// ServiceBaseOntology 基础-平台本体管理服务
	ServiceBaseOntology Service = "base-ontology"

	// 平台复用基础-身份认证服务
	ServicePlatformBaseIam Service = "platform-base-iam"
	// 平台复用用户服务
	ServicePlatformBaseUser Service = "platform-base-user"

	// ServiceSharedContent 内容服务
	ServiceSharedContent Service = "shared-content"
	// ServiceSharedFrontend 前端服务
	ServiceSharedFrontend Service = "shared-frontend"

	// ServiceSharedOperations 运营服务
	ServiceSharedOperations Service = "shared-operations"
	// ServiceSharedFacility 设施服务
	ServiceSharedFacility Service = "shared-facility"
	// ServiceSharedAI 人工智能服务
	ServiceSharedAI Service = "shared-ai"
	// ServiceSharedApproval 审批服务
	ServiceSharedApproval Service = "shared-approval"
	// ServiceSharedTask 任务服务
	ServiceSharedTask Service = "shared-task"
	// ServiceSharedResource 资源服务
	ServiceSharedResource Service = "shared-resource"
	// ServiceSharedCalendar 日程服务
	ServiceSharedCalendar Service = "shared-calendar"
	// ServiceSharedReview 审核服务
	ServiceSharedReview Service = "shared-review"
	// ServiceSharedPayment 支付服务
	ServiceSharedPayment Service = "shared-payment"
	// ServiceSharedComment 通用评论服务
	ServiceSharedComment Service = "shared-comment"

	// 平台复用内容服务
	ServicePlatformSharedContent Service = "platform-shared-content"
	// 平台复用前端服务
	ServicePlatformSharedFrontend Service = "platform-shared-frontend"
	// 平台复用运营服务
	ServicePlatformSharedOperations Service = "platform-shared-operations"

	// ServiceEduProject 教育-项目协作服务
	ServiceEduProject Service = "edu-project"
	// ServiceEduBasic 教育-基础服务
	ServiceEduBasic Service = "edu-basic"
	// ServiceEduStudent 教育-学生服务
	ServiceEduStudent Service = "edu-student"
	// ServiceEduSchedule 教育-课表服务
	ServiceEduSchedule Service = "edu-schedule"
	// ServiceEduEvaluation 教育-评价服务
	ServiceEduEvaluation Service = "edu-evaluation"
	// ServiceEduPrep 教育-备课服务
	ServiceEduPrep Service = "edu-prep"
	// ServiceEduElearning 教育-网校服务
	ServiceEduElearning Service = "edu-elearning"
	// ServiceEduBFF 教育-BFF服务
	ServiceEduBFF Service = "edu-bff"
	// 平台复用教育-基础服务
	ServicePlatformEduBasic Service = "platform-edu-basic"
)

// Values 服务列表
func ServiceValues() []string {
	return []string{
		string(ServiceBaseIam),
		string(ServiceBaseTenant),
		string(ServiceBaseNotify),
		string(ServiceBaseSafe),
		string(ServiceBaseWechat),
		string(ServiceBaseStatic),
		string(ServiceBaseStat),
		string(ServiceBaseUser),
		string(ServiceBaseRTC),
		string(ServiceBaseCRDT),
		string(ServiceBaseWorkflow),
		string(ServiceBaseAudit),
		string(ServiceBaseRag),
		string(ServiceBaseAnalysis),
		string(ServiceBaseSubscription),
		string(ServiceBaseArtifact),
		string(ServiceBaseOntology),

		string(ServicePlatformBaseIam),
		string(ServicePlatformBaseUser),

		string(ServiceSharedContent),
		string(ServiceSharedFrontend),
		string(ServiceSharedOperations),
		string(ServiceSharedFacility),
		string(ServiceSharedAI),
		string(ServiceSharedApproval),
		string(ServiceSharedTask),
		string(ServiceSharedResource),
		string(ServiceSharedCalendar),
		string(ServiceSharedReview),
		string(ServiceSharedPayment),
		string(ServiceSharedComment),

		string(ServicePlatformSharedContent),
		string(ServicePlatformSharedFrontend),
		string(ServicePlatformSharedOperations),

		string(ServiceEduProject),
		string(ServiceEduBasic),
		string(ServiceEduStudent),
		string(ServiceEduSchedule),
		string(ServiceEduEvaluation),
		string(ServiceEduPrep),
		string(ServiceEduElearning),
		string(ServiceEduBFF),

		string(ServicePlatformEduBasic),
	}
}

// ServiceLabels 返回平台可识别的服务编码与展示标签，供配置和管理界面统一枚举。
func ServiceLabels() map[string]string {
	return map[string]string{
		"base-iam":          "base-iam",
		"base-notify":       "base-notify",
		"base-safe":         "base-safe",
		"base-wechat":       "base-wechat",
		"base-static":       "base-static",
		"base-stat":         "base-stat",
		"base-user":         "base-user",
		"base-rtc":          "base-rtc",
		"base-workflow":     "base-workflow",
		"base-audit":        "base-audit",
		"base-rag":          "base-rag",
		"base-analysis":     "base-analysis",
		"base-subscription": "base-subscription",
		"base-artifact":     "base-artifact",

		"platform-base-iam":  "platform-base-iam",
		"platform-base-user": "platform-base-user",

		"shared-content":    "shared-content",
		"shared-frontend":   "shared-frontend",
		"shared-operations": "shared-operations",
		"shared-facility":   "shared-facility",
		"shared-ai":         "shared-ai",
		"shared-approval":   "shared-approval",
		"shared-task":       "shared-task",
		"shared-resource":   "shared-resource",
		"shared-calendar":   "shared-calendar",
		"shared-review":     "shared-review",
		"shared-payment":    "shared-payment",
		"shared-comment":    "shared-comment",

		"platform-shared-content":    "platform-shared-content",
		"platform-shared-frontend":   "platform-shared-frontend",
		"platform-shared-operations": "platform-shared-operations",

		"edu-project":    "edu-project",
		"edu-basic":      "edu-basic",
		"edu-student":    "edu-student",
		"edu-schedule":   "edu-schedule",
		"edu-evaluation": "edu-evaluation",
		"edu-prep":       "edu-prep",
		"edu-elearning":  "edu-elearning",
		"edu-bff":        "edu-bff",

		"platform-edu-basic": "platform-edu-basic",
	}
}

// 多语言支持
var serviceLangMap = map[string]map[string]string{
	"zh": {
		"base-iam":          "基础-身份认证",
		"base-notify":       "基础-通知",
		"base-safe":         "基础-安全",
		"base-wechat":       "基础-微信",
		"base-static":       "基础-静态",
		"base-stat":         "基础-统计",
		"base-user":         "基础-用户",
		"base-rtc":          "基础-RTC",
		"base-crdt":         "基础-CRDT",
		"base-workflow":     "基础-工作流",
		"base-audit":        "基础-审计",
		"base-subscription": "基础-订阅",

		"platform-base-iam":  "平台-基础-身份认证",
		"platform-base-user": "平台-基础-用户",

		"shared-content":    "内容服务",
		"shared-frontend":   "前端服务",
		"shared-operations": "运营服务",
		"shared-facility":   "设施服务",
		"shared-ai":         "人工智能服务",
		"shared-approval":   "审批服务",
		"shared-task":       "任务服务",
		"shared-resource":   "资源服务",
		"shared-calendar":   "日程服务",
		"shared-review":     "审核服务",
		"shared-payment":    "支付服务",
		"shared-comment":    "评论服务",

		"platform-shared-content":    "平台-内容服务",
		"platform-shared-frontend":   "平台-前端服务",
		"platform-shared-operations": "平台-运营服务",

		"edu-project":    "项目协作",
		"edu-basic":      "教育-基础",
		"edu-student":    "教育-学生",
		"edu-schedule":   "教育-课表",
		"edu-evaluation": "教育-评价",
		"edu-prep":       "教育-备课",
		"edu-elearning":  "教育-在线学习",
		"edu-bff":        "教育-BFF",

		"platform-edu-basic": "平台-教育-基础",
	},
	"ru": {
		"base-iam":          "базовое-идентификация",
		"base-notify":       "базовое-уведомление",
		"base-safe":         "базовое-безопасное",
		"base-wechat":       "базовое-вэчат",
		"base-static":       "базовое-статическое",
		"base-stat":         "базовое-статистика",
		"base-user":         "базовое-пользователь",
		"base-rtc":          "базовое-rtc",
		"base-crdt":         "базовое-crdt",
		"base-workflow":     "базовое-рабочий-поток",
		"base-audit":        "базовое-аудит",
		"base-subscription": "базовое-подписка",

		"platform-base-iam":  "платформа-базовое-идентификация",
		"platform-base-user": "платформа-базовое-пользователь",

		"shared-content":    "контент-сервис",
		"shared-frontend":   "фронтенд-сервис",
		"shared-operations": "операции-сервис",
		"shared-facility":   "инфраструктура-сервис",
		"shared-ai":         "искусственный-интеллект-сервис",
		"shared-approval":   "сервис-одобрения",
		"shared-task":       "задание-сервис",
		"shared-resource":   "ресурс-сервис",
		"shared-calendar":   "календарь-сервис",
		"shared-review":     "проверка-сервис",
		"shared-payment":    "платежный-сервис",
		"shared-comment":    "сервис-комментариев",

		"platform-shared-content":    "платформа-контент-сервис",
		"platform-shared-frontend":   "платформа-фронтенд-сервис",
		"platform-shared-operations": "платформа-операции-сервис",

		"edu-project":    "сотрудничество-проекты",
		"edu-basic":      "образование-базовое",
		"edu-student":    "образование-студент",
		"edu-schedule":   "образование-расписание",
		"edu-evaluation": "образование-оценка",
		"edu-prep":       "образование-подготовка",
		"edu-elearning":  "образование-онлайн-обучение",
		"edu-bff":        "образование-bff",

		"platform-edu-basic": "платформа-образование-базовое",
	},
}

// ServiceLangLabels _
func ServiceLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(serviceLangMap, langs)
}

// String _
func (s Service) String() string {
	return string(s)
}

// Label _
func (s Service) Label() string {
	return ServiceLabels()[s.String()]
}

// LangLabel _
func (s Service) LangLabel(langs ...string) (string, bool) {
	labels, b := ServiceLangLabels(langs...)
	if !b {
		return s.Label(), false
	}
	return labels[s.String()], true
}

// Values _
func (s Service) Values() []string {
	return ServiceValues()
}

// 服务类型
type ServiceType string

const (
	ServiceTypeBase ServiceType = "base"
	ServiceTypeEdu  ServiceType = "edu"
)

// Values _
func (s ServiceType) Values() []string {
	return []string{
		string(ServiceTypeBase),
		string(ServiceTypeEdu),
	}
}

// ServiceTypeLabels _
func ServiceTypeLabels() map[string]string {
	return map[string]string{
		"base": "base",
		"edu":  "edu",
	}
}

// 多语言支持
var serviceTypeLangMap = map[string]map[string]string{
	"zh": {
		"base": "基础",
		"edu":  "教育",
	},
	"ru": {
		"base": "базовый",
		"edu":  "образование",
	},
}

// ServiceTypeLangLabels _
func ServiceTypeLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(serviceTypeLangMap, langs)
}

// String _
func (s ServiceType) String() string {
	return string(s)
}

// Label _
func (s ServiceType) Label() string {
	return ServiceTypeLabels()[s.String()]
}

// LangLabel _
func (s ServiceType) LangLabel(langs ...string) (string, bool) {
	labels, b := ServiceTypeLangLabels(langs...)
	if !b {
		return s.Label(), false
	}
	return labels[s.String()], true
}
