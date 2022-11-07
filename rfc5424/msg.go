package rfc5424

import (
	"bytes"
	"strconv"
	"time"
)

type severity int

// rfc5424
const (
	space                      = " "
	nilValue                   = "-"
	sysLogVersion              = 1
	facility                   = 17
	UnknownSeverity   severity = -1
	EmergencySeverity severity = 0
	AlertSeverity     severity = 1
	CriticalSeverity  severity = 2
	ErrorSeverity     severity = 3
	WarningSeverity   severity = 4
	NoticeSeverity    severity = 5
	InfoSeverity      severity = 6
	DebugSeverity     severity = 7
)

func (l severity) String() string {
	switch l {
	case EmergencySeverity:
		return "EMERGENCY"
	case AlertSeverity:
		return "ALERT"
	case CriticalSeverity:
		return "CRITICAL"
	case ErrorSeverity:
		return "ERROR"
	case WarningSeverity:
		return "WARNING"
	case NoticeSeverity:
		return "NOTICE"
	case InfoSeverity:
		return "INFO"
	case DebugSeverity:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}

type syslogMsg struct {
	sev            severity
	timestamp      time.Time
	hostname       string
	appName        string
	message        string
	procId         int
	msgId          string
	structuredData map[string]map[string]string
	buf            bytes.Buffer
}

func newMsg() *syslogMsg {
	return &syslogMsg{
		sev:            UnknownSeverity,
		timestamp:      time.Time{},
		hostname:       "-",
		appName:        "-",
		message:        "-",
		procId:         -1,
		structuredData: make(map[string]map[string]string),
		buf:            bytes.Buffer{},
	}
}

func (m *syslogMsg) Severity(sev severity) *syslogMsg {
	m.sev = sev
	return m
}

func (m *syslogMsg) Timestamp(timestamp time.Time) *syslogMsg {
	m.timestamp = timestamp
	return m
}

func (m *syslogMsg) Hostname(hostname string) *syslogMsg {
	m.hostname = hostname
	return m
}

func (m *syslogMsg) AppName(appname string) *syslogMsg {
	m.appName = appname
	return m
}

func (m *syslogMsg) Message(msg string) *syslogMsg {
	m.message = msg
	return m
}

func (m *syslogMsg) MsgId(msgId string) *syslogMsg {
	m.msgId = msgId
	return m
}

func (m *syslogMsg) ProcId(procId int) *syslogMsg {
	m.procId = procId
	return m
}

func (m *syslogMsg) SdParam(sdId, name, value string) *syslogMsg {
	params, ok := m.structuredData[sdId]
	if !ok {
		m.structuredData[sdId] = map[string]string{name: value}
		return m
	}
	params[name] = value
	return m
}

func (m *syslogMsg) Build() []byte {
	m.writePriorityAndVersion()
	m.writeTimestamp()
	m.writeHostname()
	m.writeAppName()
	m.writeProcId()
	m.writMsgId()
	m.writeStructuredData()
	m.writeMsg()
	data := m.buf.Bytes()
	m.buf.Reset()
	return data
}

//	The Priority value is calculated by first multiplying the Facility
//	number by 8 and then adding the numerical value of the severity
//
// https://www.rfc-editor.org/rfc/rfc5424#section-6.2
func (m *syslogMsg) writePriorityAndVersion() {
	val := int((facility * 8) + m.sev)
	m.buf.WriteString("<" + strconv.Itoa(val) + ">" + strconv.Itoa(sysLogVersion))
}

func (m *syslogMsg) writeTimestamp() {
	if m.timestamp.IsZero() {
		m.buf.WriteString(space + nilValue)
	} else {
		m.buf.WriteString(space + m.timestamp.Format(time.RFC3339))
	}
}

func (m *syslogMsg) writeHostname() {
	if m.hostname == "" {
		m.buf.WriteString(space + nilValue)
	} else {
		m.buf.WriteString(space + m.hostname)
	}
}

func (m *syslogMsg) writeAppName() {
	if m.appName == "" {
		m.buf.WriteString(space + nilValue)
	} else {
		m.buf.WriteString(space + m.appName)
	}
}

func (m *syslogMsg) writeProcId() {
	if m.procId < 0 {
		m.buf.WriteString(space + nilValue)
	} else {
		m.buf.WriteString(space + strconv.Itoa(m.procId))
	}
}

func (m *syslogMsg) writMsgId() {
	if m.msgId == "" {
		m.buf.WriteString(space + nilValue)
	} else {
		m.buf.WriteString(space + m.msgId)
	}
}

func (m *syslogMsg) writeStructuredData() {
	if len(m.structuredData) < 1 {
		m.buf.WriteString(space + nilValue)
		return
	}
	m.buf.WriteString(space)
	for sdId, params := range m.structuredData {
		m.buf.WriteString("[" + sdId)
		for name, value := range params {
			m.buf.WriteString(space + name + "=" + strconv.Quote(value))
		}
		m.buf.WriteString("]")
	}

}

func (m *syslogMsg) writeMsg() {
	// 	bomUtf8 := []byte{0xEF, 0xBB, 0xBF}
	// TODO: If a rfc5424 application encodes MSG in UTF-8, the string MUST start
	//   with the Unicode byte order mask (BOM), which for UTF-8 is ABNF
	//   %xEF.BB.BF
	if m.message != "" {
		m.buf.WriteString(space + m.message)
	}
}
