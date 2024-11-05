package logger

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type RequestLog struct {
	ID        primitive.ObjectID  `json:"id,omitempty" bson:"_id,omitempty"`
	Referer   string              `json:"referer" bson:"referer"`
	Protocol  string              `json:"protocol" bson:"protocol"`
	Port      string              `json:"port" bson:"port"`
	IP        string              `json:"ip" bson:"ip"`
	IPs       string              `json:"ips" bson:"ips"`
	Host      string              `json:"host" bson:"host"`
	Method    string              `json:"method" bson:"method"`
	Path      string              `json:"path" bson:"path"`
	URL       string              `json:"url" bson:"url"`
	UA        string              `json:"ua" bson:"ua"`
	Latency   string              `json:"latency" bson:"latency"`
	Route     string              `json:"route" bson:"route"`
	Status    int                 `json:"status" bson:"status"`
	Headers   map[string][]string `json:"headers" bson:"headers"`
	Params    map[string]string   `json:"params" bson:"params"`
	Body      map[string]string   `json:"body" bson:"body"`
	Form      map[string]string   `json:"form" bson:"form"`
	Cookies   map[string]string   `json:"cookies" bson:"cookies"`
	CreatedAt time.Time           `json:"createdAt" bson:"createdAt"`
}

func NewLog(c *fiber.Ctx) *RequestLog {
	return &RequestLog{
		Referer:   c.Get(fiber.HeaderReferer),
		Protocol:  c.Protocol(),
		Port:      c.Port(),
		IP:        c.IP(),
		IPs:       c.Get(fiber.HeaderXForwardedFor),
		Host:      c.Hostname(),
		Path:      c.Path(),
		URL:       c.OriginalURL(),
		UA:        c.Get(fiber.HeaderUserAgent),
		Route:     c.Route().Path,
		Body:      make(map[string]string),
		Headers:   c.GetReqHeaders(),
		Params:    c.Queries(),
		Status:    c.Response().StatusCode(),
		Method:    c.Method(),
		Cookies:   make(map[string]string),
		Form:      make(map[string]string),
		Latency:   "0",
		CreatedAt: time.Now(),
	}
}
