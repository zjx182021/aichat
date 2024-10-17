package metricsbus

import "github.com/prometheus/client_golang/prometheus"

type BusMetrics struct {
	QuestionsTotalCounter          prometheus.Counter
	KeywordsQuestionsTotalCounter  prometheus.Counter
	SensitiveQuestionsTotalCounter prometheus.Counter
	ErrQuestionsTotalCounter       prometheus.Counter
}

const (
	NAMESPACE = "ai_chat"
	SUBSYSTEM = "chat_service"
)

func NewBusMetrics(registry *prometheus.Registry) *BusMetrics {
	questionsTotalCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   NAMESPACE,
		Subsystem:   SUBSYSTEM,
		Name:        "questions_total",
		ConstLabels: map[string]string{"app": "ai_chat"},
		Help:        "记录用户提交问题的总数，仅包含记录到DB的问题数量",
	})
	keywordsQuestionsTotalCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   NAMESPACE,
		Subsystem:   SUBSYSTEM,
		Name:        "keywords_questions_total",
		ConstLabels: map[string]string{"app": "ai_chat"},
		Help:        "记录用户提交的包含关键词的问题总数",
	})
	sensitiveQuestionsTotalCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   NAMESPACE,
		Subsystem:   SUBSYSTEM,
		Name:        "sensitive_questions_total",
		ConstLabels: map[string]string{"app": "ai_chat"},
		Help:        "记录用户提交的触发敏感词的问题总数",
	})
	errQuestionsTotalCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   NAMESPACE,
		Subsystem:   SUBSYSTEM,
		Name:        "err_questions_total",
		ConstLabels: map[string]string{"app": "ai_chat"},
		Help:        "记录用户提交问题时报错的总数",
	})
	registry.MustRegister(questionsTotalCounter, keywordsQuestionsTotalCounter, sensitiveQuestionsTotalCounter, errQuestionsTotalCounter)
	return &BusMetrics{
		QuestionsTotalCounter:          questionsTotalCounter,
		KeywordsQuestionsTotalCounter:  keywordsQuestionsTotalCounter,
		SensitiveQuestionsTotalCounter: sensitiveQuestionsTotalCounter,
		ErrQuestionsTotalCounter:       errQuestionsTotalCounter,
	}
}
