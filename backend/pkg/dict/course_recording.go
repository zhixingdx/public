package dict

import "github.com/zhixingdx/oipivot/backend/pkg/lang"

type CourseRecordingSource string

const (
	// 本地上传
	CourseRecordingSourceLocalUpload CourseRecordingSource = "local_upload"
	// 录制转存
	CourseRecordingSourceRecordSaved CourseRecordingSource = "record_saved"
	// 资源转存
	CourseRecordingSourceResourceSaved CourseRecordingSource = "resource_saved"
)

func CourseRecordingSourceValues() []string {
	return []string{
		string(CourseRecordingSourceLocalUpload),
		string(CourseRecordingSourceRecordSaved),
		string(CourseRecordingSourceResourceSaved),
	}
}

func CourseRecordingSourceLabels() map[string]string {
	return map[string]string{
		"local_upload":   "本地上传",
		"record_saved":   "录制转存",
		"resource_saved": "资源转存",
	}
}

func CourseRecordingSourceLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(map[string]map[string]string{
		"zh": {
			"local_upload":   "本地上传",
			"record_saved":   "录制转存",
			"resource_saved": "资源转存",
		},
		"ru": {
			"local_upload":   "локальный загрузка",
			"record_saved":   "запись сохранена",
			"resource_saved": "ресурс сохранен",
		},
	}, langs)
}

func (s CourseRecordingSource) String() string {
	return string(s)
}

func (s CourseRecordingSource) Label() string {
	return CourseRecordingSourceLabels()[s.String()]
}

func (s CourseRecordingSource) LangLabel(langs ...string) (string, bool) {
	labels, b := CourseRecordingSourceLangLabels(langs...)
	if !b {
		return s.Label(), false
	}
	return labels[s.String()], true
}

func (s CourseRecordingSource) Values() []string {
	return CourseRecordingSourceValues()
}
