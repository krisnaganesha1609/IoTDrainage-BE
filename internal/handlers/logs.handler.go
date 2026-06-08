package handlers

import (
	"bytes"
	"html/template"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/entities"
	"github.com/krisnaganesha1609/IoTDrainage-BE/utils"
)

// LogsPageData is passed to the HTML template.
type LogsPageData struct {
	Logs []entities.SensorLogEntity
}

const logsHTMLTemplate = `<!DOCTYPE html>
<html lang="id">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>System Logs — IoT Drainage Monitor</title>
<style>
  :root { --bg:#0f1117; --surface:#1a1d27; --border:#2d3048; --text:#e2e8f0; --muted:#8892a4;
    --info:#38bdf8; --warn:#fbbf24; --err:#f87171; --green:#4ade80; }
  *{box-sizing:border-box;margin:0;padding:0;}
  body{background:var(--bg);color:var(--text);font-family:'Segoe UI',system-ui,sans-serif;font-size:14px;}
  header{background:var(--surface);border-bottom:1px solid var(--border);padding:18px 24px;display:flex;align-items:center;gap:12px;}
  header h1{font-size:18px;font-weight:600;}
  header span{font-size:12px;color:var(--muted);margin-left:auto;}
  .container{max-width:1280px;margin:0 auto;padding:24px;}
  .filter-bar{display:flex;gap:10px;margin-bottom:20px;flex-wrap:wrap;align-items:center;}
  .filter-bar input{background:var(--surface);border:1px solid var(--border);color:var(--text);padding:8px 12px;border-radius:6px;font-size:13px;width:280px;outline:none;}
  .filter-bar input:focus{border-color:var(--info);}
  .legend{display:flex;gap:16px;margin-left:auto;font-size:12px;}
  .dot{width:8px;height:8px;border-radius:50%;display:inline-block;margin-right:5px;}
  .dot-info{background:var(--info);}.dot-warn{background:var(--warn);}.dot-err{background:var(--err);}
  table{width:100%;border-collapse:collapse;background:var(--surface);border-radius:8px;overflow:hidden;}
  thead{background:#12141e;}
  th{padding:10px 14px;text-align:left;font-size:11px;text-transform:uppercase;letter-spacing:.06em;color:var(--muted);border-bottom:1px solid var(--border);}
  td{padding:10px 14px;border-bottom:1px solid var(--border);vertical-align:top;font-size:13px;}
  tr:last-child td{border-bottom:none;}
  tr:hover td{background:rgba(255,255,255,.025);}
  .badge{display:inline-block;padding:2px 8px;border-radius:4px;font-size:11px;font-weight:600;letter-spacing:.04em;}
  .level-info .badge{background:rgba(56,189,248,.15);color:var(--info);}
  .level-warning .badge{background:rgba(251,191,36,.15);color:var(--warn);}
  .level-error .badge{background:rgba(248,113,113,.15);color:var(--err);}
  .mono{font-family:'Cascadia Code','Fira Code',monospace;font-size:12px;}
  .empty{text-align:center;padding:48px;color:var(--muted);}
  .summary{display:flex;gap:12px;margin-bottom:20px;flex-wrap:wrap;}
  .card{background:var(--surface);border:1px solid var(--border);border-radius:8px;padding:14px 18px;flex:1;min-width:120px;}
  .card-label{font-size:11px;color:var(--muted);text-transform:uppercase;letter-spacing:.06em;}
  .card-value{font-size:28px;font-weight:700;margin-top:4px;}
  .card-info .card-value{color:var(--info);}.card-warn .card-value{color:var(--warn);}
  .card-err .card-value{color:var(--err);}.card-total .card-value{color:var(--green);}
  .msg{max-width:280px;word-break:break-word;}
</style>
</head>
<body>
<header>
  <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#38bdf8" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>
  <h1>System Logs — IoT Drainage Monitor</h1>
  <span>7 hari terakhir &bull; refresh halaman untuk data terbaru</span>
</header>
<div class="container">
  <div class="summary">
    <div class="card card-total"><div class="card-label">Total Logs</div><div class="card-value">{{len .Logs}}</div></div>
    <div class="card card-info"><div class="card-label">INFO</div><div class="card-value">{{.CountInfo}}</div></div>
    <div class="card card-warn"><div class="card-label">WARNING</div><div class="card-value">{{.CountWarn}}</div></div>
    <div class="card card-err"><div class="card-label">ERROR</div><div class="card-value">{{.CountErr}}</div></div>
  </div>
  <div class="filter-bar">
    <input id="search" type="text" placeholder="🔍 Cari pesan, device ID, level..." oninput="filterLogs()">
    <div class="legend">
      <span><span class="dot dot-info"></span>INFO</span>
      <span><span class="dot dot-warn"></span>WARNING</span>
      <span><span class="dot dot-err"></span>ERROR</span>
    </div>
  </div>
  {{if eq (len .Logs) 0}}
  <div class="empty">Belum ada log dalam 7 hari terakhir.</div>
  {{else}}
  <table id="logTable">
    <thead><tr>
      <th>Waktu (WIB)</th><th>Device ID</th><th>Level</th><th>Pesan</th>
      <th>Wake / Reset</th><th>Heap (byte)</th><th>RSSI (dBm)</th>
      <th>Net Fail</th><th>Aktif (ms)</th>
    </tr></thead>
    <tbody>
    {{range .Logs}}
    <tr class="{{levelClass .Level}}">
      <td class="mono">{{formatTime .Time}}</td>
      <td class="mono">{{.DeviceID}}</td>
      <td><span class="badge">{{.Level}}</span></td>
      <td class="msg">{{.Message}}</td>
      <td class="mono">{{.WakeReason}} / {{.ResetReason}}</td>
      <td class="mono">{{.FreeHeapBytes}}</td>
      <td class="mono">{{.WifiRssiDbm}}</td>
      <td class="mono">{{.NetworkFailures}}</td>
      <td class="mono">{{.ActiveTimeMs}}</td>
    </tr>
    {{end}}
    </tbody>
  </table>
  {{end}}
</div>
<script>
function filterLogs(){
  const q=document.getElementById('search').value.toLowerCase();
  document.querySelectorAll('#logTable tbody tr').forEach(r=>{
    r.style.display=r.textContent.toLowerCase().includes(q)?'':'none';
  });
}
</script>
</body></html>`

// ExtendedLogsPageData adds pre-counted level stats so the template stays simple.
type ExtendedLogsPageData struct {
	Logs      []entities.SensorLogEntity
	CountInfo int
	CountWarn int
	CountErr  int
}

// ShowLogsPage renders GET /logs as a server-side HTML page showing
// device system logs from the last 7 days (optionally filtered by ?device_id=).
func (h *Handler) ShowLogsPage(c fiber.Ctx) error {
	deviceID := c.Query("device_id")

	logs, err := h.Service.GetRecentLogs(deviceID)
	if err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}

	data := ExtendedLogsPageData{Logs: logs}
	for _, l := range logs {
		switch l.Level {
		case "INFO":
			data.CountInfo++
		case "WARNING":
			data.CountWarn++
		case "ERROR":
			data.CountErr++
		}
	}

	funcMap := template.FuncMap{
		"formatTime": func(t time.Time) string {
			return t.In(time.FixedZone("WIB", 7*3600)).Format("02 Jan 2006  15:04:05 WIB")
		},
		"levelClass": func(level string) string {
			switch level {
			case "ERROR":
				return "level-error"
			case "WARNING":
				return "level-warning"
			default:
				return "level-info"
			}
		},
	}

	tmpl, parseErr := template.New("logs").Funcs(funcMap).Parse(logsHTMLTemplate)
	if parseErr != nil {
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Template parse error: "+parseErr.Error())
	}

	var buf bytes.Buffer
	if execErr := tmpl.Execute(&buf, data); execErr != nil {
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Template render error: "+execErr.Error())
	}

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.Status(fiber.StatusOK).Send(buf.Bytes())
}
