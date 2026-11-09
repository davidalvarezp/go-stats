package server

import (
	"strconv"
	"strings"
)

// GetDashboardHTML returns the complete, standalone single-page dashboard.
// All styles, scripts, and layouts are self-contained with ZERO external node/npm
// dependencies, ensuring that `go build` produces a 100% monolithic, portable binary.
func GetDashboardHTML(port int) string {
	return strings.ReplaceAll(dashboardTemplate, "{{PORT}}", strconv.Itoa(port))
}

const dashboardTemplate = `<!DOCTYPE html>
<html lang="en" class="dark">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>go-stats · Linux System Monitor</title>
  <meta name="description" content="Lightweight native Linux system monitor written in Go, running as a systemd service.">
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;600;700&family=Plus+Jakarta+Sans:wght@400;500;600;700&display=swap" rel="stylesheet">
  <style>
    /* ==========================================================================
       BRANDING & THEME COLOR VARIABLES (CUSTOMIZE HERE)
       --------------------------------------------------------------------------
       Author: davidalvarezp (dap)
       Website: https://davidalvarezp.com/go-stats
       Repository: https://github.com/davidalvarezp/go-stats
       ========================================================================== */
    :root {
      /* [Dark Theme - Default Linux Terminal Green & Deep Charcoal] */
      --bg-canvas: #090a0c;             /* Main application background (deep terminal black) */
      --bg-surface: #121418;            /* Primary cards, containers, and table backgrounds */
      --bg-surface-elevated: #1a1d24;   /* Hover states, inputs, and selected items */
      --border-subtle: #242933;         /* Hairline dividers, card outlines, table borders */
      --border-focus: #16a34a;          /* Active focus rings and highlighted borders */
      
      --text-primary: #f0f3f6;          /* High-contrast headings and primary values */
      --text-secondary: #8b949e;        /* Metadata labels, column headers, units */
      --text-muted: #57606a;            /* Inactive tabs, disabled states, subtle hints */

      /* Accent Colors */
      --color-accent-green: #22c55e;    /* Signature terminal green: logo, healthy status, CPU */
      --color-accent-glow: rgba(34, 197, 94, 0.15); /* Green translucent badges and progress bars */
      --color-cyan: #06b6d4;            /* Secondary telemetry: disk read/write, virtual mem */
      --color-indigo: #818cf8;          /* Swap memory and buffered cache metrics */
      --color-warning: #eab308;         /* Elevated CPU (>60%), disk wait, stopped processes */
      --color-danger: #ef4444;          /* High CPU (>85%), zombie processes, service errors */
      
      /* Typography */
      --font-sans: 'Plus Jakarta Sans', system-ui, -apple-system, sans-serif;
      --font-mono: 'JetBrains Mono', ui-monospace, SFMono-Regular, monospace;
    }

    /* [Light Theme - Clean Engineering Contrast] */
    html.light {
      --bg-canvas: #f8fafc;             /* Crisp light background */
      --bg-surface: #ffffff;            /* Pure white card backgrounds */
      --bg-surface-elevated: #f1f5f9;   /* Light gray hover backgrounds */
      --border-subtle: #e2e8f0;         /* Subtle gray border */
      --border-focus: #15803d;          /* Darker green focus ring for light background */
      
      --text-primary: #0f172a;          /* Deep slate black text */
      --text-secondary: #475569;        /* Medium slate text */
      --text-muted: #94a3b8;            /* Soft placeholder text */

      --color-accent-green: #16a34a;    /* High-contrast emerald green */
      --color-accent-glow: rgba(22, 163, 74, 0.12);
      --color-cyan: #0891b2;
      --color-indigo: #6366f1;
      --color-warning: #d97706;
      --color-danger: #dc2626;
    }

    /* Base Reset */
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background: var(--bg-canvas);
      color: var(--text-primary);
      font-family: var(--font-sans);
      font-size: 13px;
      line-height: 1.5;
      -webkit-font-smoothing: antialiased;
      transition: background 0.15s ease, color 0.15s ease;
    }
    .mono { font-family: var(--font-mono); }
    .container { max-width: 1440px; margin: 0 auto; padding: 24px; }
    a { color: var(--color-accent-green); text-decoration: none; }
    a:hover { text-decoration: underline; }

    /* Top Navigation Header */
    header {
      background: var(--bg-surface);
      border-bottom: 1px solid var(--border-subtle);
      position: sticky;
      top: 0;
      z-index: 50;
    }
    .header-inner {
      max-width: 1440px;
      margin: 0 auto;
      padding: 12px 24px;
      display: flex;
      justify-content: space-between;
      align-items: center;
      flex-wrap: wrap;
      gap: 16px;
    }
    .brand-group { display: flex; align-items: center; gap: 12px; }
    .brand-badge {
      width: 38px; height: 38px; border-radius: 8px;
      background: var(--color-accent-glow);
      border: 1px solid rgba(34, 197, 94, 0.3);
      color: var(--color-accent-green);
      display: flex; align-items: center; justify-content: center;
      font-weight: 700; font-size: 15px; font-family: var(--font-mono);
    }
    .brand-title { font-size: 16px; font-weight: 700; letter-spacing: -0.3px; }
    .brand-meta { font-size: 11px; color: var(--text-secondary); display: flex; align-items: center; gap: 6px; }

    /* Badges & Tags */
    .badge {
      font-size: 11px; padding: 3px 8px; border-radius: 6px; font-weight: 600;
      background: var(--bg-surface-elevated); color: var(--text-secondary);
      border: 1px solid var(--border-subtle); display: inline-flex; align-items: center; gap: 6px;
    }
    .badge-green {
      background: var(--color-accent-glow); color: var(--color-accent-green);
      border-color: rgba(34, 197, 94, 0.3);
    }
    .badge-danger {
      background: rgba(239, 68, 68, 0.15); color: var(--color-danger);
      border-color: rgba(239, 68, 68, 0.3);
    }
    .dot { width: 6px; height: 6px; border-radius: 50%%; background: currentColor; }
    .dot-pulse { animation: pulse 2s infinite ease-in-out; }
    @keyframes pulse { 0%%, 100%% { opacity: 1; } 50%% { opacity: 0.4; } }

    /* Button Controls */
    .btn-group { display: flex; align-items: center; gap: 8px; }
    button {
      background: var(--bg-surface-elevated); color: var(--text-primary);
      border: 1px solid var(--border-subtle); padding: 6px 12px; border-radius: 6px;
      cursor: pointer; font-size: 12px; font-weight: 500; font-family: inherit;
      display: inline-flex; align-items: center; gap: 6px; transition: all 0.15s ease;
    }
    button:hover { border-color: var(--text-secondary); }
    button.primary {
      background: var(--color-accent-green); color: #000;
      border-color: var(--color-accent-green); font-weight: 600;
    }
    button.primary:hover { opacity: 0.9; }
    button.danger {
      background: rgba(239, 68, 68, 0.12); color: var(--color-danger);
      border-color: rgba(239, 68, 68, 0.3);
    }
    button.danger:hover { background: rgba(239, 68, 68, 0.2); }

    /* Navigation Tabs */
    .tabs-bar {
      display: flex; gap: 4px; padding: 0 24px; max-width: 1440px; margin: 0 auto;
      border-bottom: 1px solid var(--border-subtle); overflow-x: auto;
    }
    .nav-tab {
      padding: 10px 16px; border: none; border-bottom: 2px solid transparent;
      background: transparent; color: var(--text-secondary); font-size: 13px;
      font-weight: 500; border-radius: 0; cursor: pointer; white-space: nowrap;
    }
    .nav-tab:hover { color: var(--text-primary); }
    .nav-tab.active {
      color: var(--color-accent-green); border-bottom-color: var(--color-accent-green); font-weight: 600;
    }

    /* Structured Cards */
    .card {
      background: var(--bg-surface); border: 1px solid var(--border-subtle);
      border-radius: 12px; padding: 20px; margin-bottom: 20px;
    }
    .card-header {
      display: flex; justify-content: space-between; align-items: center;
      margin-bottom: 16px; flex-wrap: wrap; gap: 8px;
    }
    .card-title { font-size: 14px; font-weight: 600; letter-spacing: -0.2px; }

    /* Classic top-style System Summary Banner */
    .top-summary {
      background: var(--bg-surface); border: 1px solid var(--border-subtle);
      border-radius: 12px; padding: 18px 22px; margin-bottom: 24px; font-size: 12px;
    }
    .top-header-line {
      display: flex; flex-wrap: wrap; gap: 8px 16px; align-items: center;
      padding-bottom: 14px; margin-bottom: 14px; border-bottom: 1px solid var(--border-subtle);
    }
    .top-columns { display: grid; grid-template-columns: 1fr 1fr; gap: 24px; }
    @media(max-width: 960px) { .top-columns { grid-template-columns: 1fr; } }
    .metric-row { margin-bottom: 10px; }

    /* Progress & Metric Bars */
    .bar-track {
      background: rgba(255, 255, 255, 0.06); height: 7px;
      border-radius: 99px; overflow: hidden; margin-top: 6px; display: flex;
    }
    html.light .bar-track { background: #e2e8f0; }
    .bar-fill { height: 100%%; transition: width 0.3s ease; }

    /* CPU Multi-Core Matrix */
    .cpu-matrix { display: grid; grid-template-columns: repeat(auto-fill, minmax(210px, 1fr)); gap: 12px; }
    .core-cell {
      background: var(--bg-surface-elevated); border: 1px solid var(--border-subtle);
      border-radius: 8px; padding: 12px;
    }
    .sparkline-container {
      height: 60px; display: flex; align-items: flex-end; gap: 3px;
      background: rgba(0, 0, 0, 0.25); padding: 8px; border-radius: 8px;
      border: 1px solid var(--border-subtle); margin-bottom: 16px;
    }
    .spark-segment {
      flex: 1; border-radius: 2px 2px 0 0; background: var(--color-accent-green);
      min-height: 4px; transition: height 0.2s ease;
    }

    /* Process Table (htop view) */
    .table-responsive { overflow-x: auto; max-height: 600px; }
    table { width: 100%%; border-collapse: collapse; text-align: left; font-size: 12px; }
    th {
      position: sticky; top: 0; background: var(--bg-surface);
      color: var(--text-secondary); padding: 10px 14px;
      border-bottom: 1px solid var(--border-subtle); font-weight: 600; cursor: pointer;
      user-select: none;
    }
    th:hover { color: var(--text-primary); }
    td { padding: 9px 14px; border-bottom: 1px solid rgba(255, 255, 255, 0.03); }
    html.light td { border-bottom: 1px solid rgba(0, 0, 0, 0.04); }
    tr:hover { background: rgba(34, 197, 94, 0.04); }
    
    .proc-state { padding: 2px 6px; border-radius: 4px; font-size: 10px; font-weight: 700; }
    .state-R { background: rgba(34, 197, 94, 0.15); color: var(--color-accent-green); }
    .state-S { background: rgba(6, 182, 212, 0.15); color: var(--color-cyan); }
    .state-Z { background: rgba(239, 68, 68, 0.15); color: var(--color-danger); }
    .state-D { background: rgba(234, 179, 8, 0.15); color: var(--color-warning); }

    /* Filter Controls */
    .filter-wrapper { display: flex; gap: 10px; flex-wrap: wrap; margin-bottom: 16px; }
    input[type="text"] {
      background: var(--bg-canvas); border: 1px solid var(--border-subtle);
      color: var(--text-primary); padding: 8px 12px; border-radius: 6px;
      font-size: 12px; font-family: inherit; outline: none; min-width: 280px;
    }
    input[type="text"]:focus { border-color: var(--border-focus); }

    /* Log Stream Terminal */
    .log-stream {
      background: #050507; color: #d4d4d8; font-family: var(--font-mono);
      padding: 16px; border-radius: 8px; border: 1px solid var(--border-subtle);
      max-height: 380px; overflow-y: auto; font-size: 11px; line-height: 1.6;
    }
    .log-row { margin-bottom: 4px; word-break: break-all; }

    /* Footer */
    footer {
      border-top: 1px solid var(--border-subtle); padding: 18px 24px;
      color: var(--text-secondary); font-size: 12px; margin-top: 40px;
    }
    .footer-inner {
      max-width: 1440px; margin: 0 auto; display: flex;
      justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 12px;
    }
  </style>
</head>
<body>

  <!-- Top Navigation Header -->
  <header>
    <div class="header-inner">
      <div class="brand-group">
        <div class="brand-badge">gs</div>
        <div>
          <div style="display:flex; align-items:center; gap:8px;">
            <span class="brand-title mono">go-stats</span>
            <span class="badge mono badge-green">:{{PORT}}</span>
            <span class="badge mono">v1.0.0</span>
          </div>
          <div class="brand-meta">
            <span>By <a href="https://davidalvarezp.com" target="_blank" rel="noopener">davidalvarezp (dap)</a></span>
            <span>·</span>
            <span><a href="https://github.com/davidalvarezp/go-stats" target="_blank" rel="noopener">GitHub</a></span>
            <span>·</span>
            <span><a href="https://davidalvarezp.com/go-stats" target="_blank" rel="noopener">Docs</a></span>
          </div>
        </div>
      </div>

      <div style="display:flex; align-items:center; gap:12px; flex-wrap:wrap;">
        <div id="serviceBadge" class="badge badge-green mono" onclick="selectTab('systemd')" title="Click to manage systemd service" style="cursor:pointer;">
          <span class="dot dot-pulse"></span>
          <span id="serviceStatusLabel">go-stats.service (active)</span>
        </div>

        <div class="btn-group">
          <button id="pauseToggleBtn" onclick="togglePause()">Pause</button>
          <button onclick="toggleTheme()" title="Toggle Dark/Light Mode">Theme</button>
        </div>
      </div>
    </div>

    <!-- Navigation Tabs -->
    <div class="tabs-bar">
      <button class="nav-tab active" onclick="selectTab('overview')">Overview (top)</button>
      <button class="nav-tab" onclick="selectTab('processes')">Processes (htop)</button>
      <button class="nav-tab" onclick="selectTab('systemd')">Systemd Service</button>
      <button class="nav-tab" onclick="selectTab('proc')">Kernel /proc</button>
    </div>
  </header>

  <main class="container">

    <!-- TAB 1: Overview (top metrics) -->
    <div id="tab-overview">
      <div class="top-summary mono">
        <div class="top-header-line">
          <span style="font-weight:700; color:var(--color-accent-green);">top -</span>
          <span id="summaryClock">--:--:--</span>
          <span>·</span>
          <span>up <strong id="summaryUptime">--</strong></span>
          <span>·</span>
          <span>1 user</span>
          <span>·</span>
          <span>load average: <strong id="summaryLoad" style="color:var(--color-accent-green)">0.00, 0.00, 0.00</strong></span>
        </div>

        <div class="top-columns">
          <div>
            <div class="metric-row">
              <strong>Tasks:</strong> <span id="metricTasksTotal">0</span> total,
              <span style="color:var(--color-accent-green); font-weight:600;"><span class="dot" style="display:inline-block; margin-right:4px;"></span><span id="metricTasksRun">0</span> running</span>,
              <span id="metricTasksSleep">0</span> sleeping,
              <span id="metricTasksStop">0</span> stopped,
              <span style="color:var(--color-danger)"><span id="metricTasksZombie">0</span> zombie</span>
            </div>

            <div class="metric-row" style="color:var(--text-secondary);">
              <strong>%%Cpu(s):</strong>
              <span style="color:var(--color-accent-green)"><strong id="cpuUserVal">0.0</strong> us</span>,
              <span style="color:var(--color-indigo)"><strong id="cpuSysVal">0.0</strong> sy</span>,
              <span>0.0 ni</span>,
              <span style="color:var(--color-accent-green)"><strong id="cpuIdleVal">100.0</strong> id</span>,
              <span style="color:var(--color-warning)"><strong id="cpuWaitVal">0.0</strong> wa</span>
            </div>

            <div>
              <div style="display:flex; justify-content:space-between; font-size:11px; margin-top:8px;">
                <span>Total CPU Load</span>
                <strong id="cpuTotalText" style="color:var(--color-accent-green)">0.0%%</strong>
              </div>
              <div class="bar-track">
                <div id="cpuTotalBar" class="bar-fill" style="background:var(--color-accent-green); width:0%%"></div>
              </div>
            </div>
          </div>

          <div>
            <div class="metric-row">
              <strong>MiB Mem:</strong> <span id="memTotalTxt">0</span> total, <span id="memFreeTxt">0</span> free, <strong id="memUsedTxt" style="color:var(--color-accent-green)">0</strong> used, <span id="memCacheTxt">0</span> buff/cache
            </div>
            <div class="bar-track">
              <div id="memUsedBar" class="bar-fill" style="background:var(--color-accent-green); width:0%%"></div>
              <div id="memCacheBar" class="bar-fill" style="background:var(--color-indigo); width:0%%"></div>
            </div>

            <div class="metric-row" style="margin-top:14px;">
              <strong>MiB Swap:</strong> <span id="swapTotalTxt">0</span> total, <span id="swapFreeTxt">0</span> free, <strong id="swapUsedTxt" style="color:var(--color-indigo)">0</strong> used
            </div>
            <div class="bar-track">
              <div id="swapUsedBar" class="bar-fill" style="background:var(--color-indigo); width:0%%"></div>
            </div>
          </div>
        </div>
      </div>

      <!-- CPU Core Matrix -->
      <div class="card">
        <div class="card-header">
          <span class="card-title">Multi-Core CPU Matrix</span>
          <span class="mono" style="font-size:12px; color:var(--text-secondary)">Activity history (last 30s)</span>
        </div>
        <div id="sparklineGrid" class="sparkline-container"></div>
        <div id="coresContainer" class="cpu-matrix mono"></div>
      </div>

      <!-- Storage and Network Grid -->
      <div style="display:grid; grid-template-columns: repeat(auto-fit, minmax(320px, 1fr)); gap:20px;">
        <div class="card mono">
          <div class="card-header">
            <span class="card-title">Storage Mounts (/proc/mounts)</span>
          </div>
          <div id="mountsContainer"></div>
        </div>

        <div class="card mono">
          <div class="card-header">
            <span class="card-title">Network Interfaces (/proc/net/dev)</span>
          </div>
          <div id="netContainer"></div>
        </div>
      </div>
    </div>

    <!-- TAB 2: Processes (htop view) -->
    <div id="tab-processes" style="display:none;">
      <div class="card">
        <div class="card-header">
          <span class="card-title">Process Activity (htop/top inspect)</span>
          <span id="processCountLabel" class="badge mono">0 processes</span>
        </div>

        <div class="filter-wrapper">
          <input type="text" id="searchInput" placeholder="Filter by command, PID, name, or user..." oninput="applyFilters()">
          <button onclick="setFilterState('ALL')" id="btn-ALL" class="primary">All</button>
          <button onclick="setFilterState('R')" id="btn-R">Running</button>
          <button onclick="setFilterState('S')" id="btn-S">Sleeping</button>
          <button onclick="setFilterState('Z')" id="btn-Z">Zombie</button>
        </div>

        <div class="table-responsive">
          <table class="mono">
            <thead>
              <tr>
                <th onclick="toggleSort('pid')">PID ↕</th>
                <th onclick="toggleSort('user')">USER ↕</th>
                <th>PR</th>
                <th>NI</th>
                <th onclick="toggleSort('virt_bytes')">VIRT ↕</th>
                <th onclick="toggleSort('res_bytes')">RES ↕</th>
                <th onclick="toggleSort('state')">S</th>
                <th onclick="toggleSort('cpu_percent')">%%CPU ↕</th>
                <th onclick="toggleSort('mem_percent')">%%MEM ↕</th>
                <th onclick="toggleSort('time_seconds')">TIME+ ↕</th>
                <th onclick="toggleSort('name')">COMMAND ↕</th>
                <th>SIGNAL</th>
              </tr>
            </thead>
            <tbody id="processesTbody">
              <tr><td colspan="12" style="text-align:center; color:var(--text-secondary); padding:24px;">Scanning /proc processes...</td></tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- TAB 3: Systemd Service Manager -->
    <div id="tab-systemd" style="display:none;">
      <div class="card">
        <div class="card-header">
          <span class="card-title">Systemd Daemon Management (go-stats.service)</span>
          <span class="badge mono">systemctl</span>
        </div>
        <p style="color:var(--text-secondary); margin-bottom:18px;">
          Controls the daemon running in the background on your Debian Linux host.
        </p>

        <div style="display:flex; gap:10px; flex-wrap:wrap; margin-bottom:24px;">
          <button class="primary" onclick="execSystemd('start')">sudo systemctl start go-stats</button>
          <button class="danger" onclick="execSystemd('stop')">sudo systemctl stop go-stats</button>
          <button onclick="execSystemd('restart')">sudo systemctl restart go-stats</button>
          <button onclick="execSystemd('enable')">sudo systemctl enable go-stats</button>
        </div>

        <div class="card-header">
          <span class="card-title">journalctl -u go-stats -f (Live logs)</span>
        </div>
        <div id="journalOutput" class="log-stream"></div>
      </div>
    </div>

    <!-- TAB 4: Kernel /proc Raw -->
    <div id="tab-proc" style="display:none;">
      <div class="card mono">
        <div class="card-header">
          <span class="card-title">Raw Linux Kernel Telemetry</span>
          <div style="display:flex; gap:6px;">
            <button onclick="fetchRaw('/proc/stat')">/proc/stat</button>
            <button onclick="fetchRaw('/proc/meminfo')">/proc/meminfo</button>
            <button onclick="fetchRaw('/proc/loadavg')">/proc/loadavg</button>
            <button onclick="fetchRaw('/proc/uptime')">/proc/uptime</button>
          </div>
        </div>
        <pre id="rawBox" class="log-stream" style="color:var(--color-accent-green); max-height:480px;"></pre>
      </div>
    </div>

  </main>

  <!-- Clean Footer with author links -->
  <footer>
    <div class="footer-inner mono">
      <div>
        <strong>go-stats</strong> · Native Linux System Monitor · Developed by <a href="https://davidalvarezp.com" target="_blank" rel="noopener">davidalvarezp (dap)</a>
      </div>
      <div>
        <a href="https://davidalvarezp.com/go-stats" target="_blank" rel="noopener">Documentation</a> · 
        <a href="https://github.com/davidalvarezp/go-stats" target="_blank" rel="noopener">GitHub</a> · 
        <a href="mailto:mail@davidalvarezp.com">mail@davidalvarezp.com</a>
      </div>
    </div>
  </footer>

  <script>
    let isStreamPaused = false;
    let sparklineHistory = Array(30).fill(12);
    let processRecords = [];
    let stateFilterMode = 'ALL';
    let sortColumnKey = 'cpu_percent';
    let isSortAscending = false;

    function formatByteUnits(bytes) {
      if (!bytes || bytes === 0) return '0 B';
      const k = 1024, dm = 1;
      const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
      const i = Math.floor(Math.log(bytes) / Math.log(k));
      return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
    }

    function toggleTheme() {
      document.documentElement.classList.toggle('light');
    }

    function togglePause() {
      isStreamPaused = !isStreamPaused;
      document.getElementById('pauseToggleBtn').innerText = isStreamPaused ? 'Resume' : 'Pause';
    }

    function selectTab(tabKey) {
      ['overview', 'processes', 'systemd', 'proc'].forEach(t => {
        document.getElementById('tab-' + t).style.display = (t === tabKey) ? 'block' : 'none';
      });
      document.querySelectorAll('.nav-tab').forEach((el, idx) => {
        el.classList.toggle('active', ['overview', 'processes', 'systemd', 'proc'][idx] === tabKey);
      });
      if (tabKey === 'proc') fetchRaw('/proc/stat');
      if (tabKey === 'systemd') fetchJournalLogs();
    }

    async function pollTelemetry() {
      if (isStreamPaused) return;
      try {
        const [mRes, pRes] = await Promise.all([fetch('/api/metrics'), fetch('/api/processes')]);
        if (mRes.ok) {
          const metrics = await mRes.json();
          displayMetrics(metrics);
        }
        if (pRes.ok) {
          const procs = await pRes.json();
          processRecords = procs.processes || [];
          document.getElementById('processCountLabel').innerText = processRecords.length + ' processes';
          applyFilters();
        }
      } catch (err) {
        console.error('Polling error', err);
      }
    }

    function displayMetrics(m) {
      document.getElementById('summaryClock').innerText = new Date(m.timestamp).toLocaleTimeString();
      document.getElementById('summaryUptime').innerText = m.uptime_string || '--';
      document.getElementById('summaryLoad').innerText = (m.load.load1 || 0).toFixed(2) + ', ' + (m.load.load5 || 0).toFixed(2) + ', ' + (m.load.load15 || 0).toFixed(2);

      const s = m.task_summary || { total: processRecords.length, running: 1, sleeping: processRecords.length - 1, stopped: 0, zombie: 0 };
      document.getElementById('metricTasksTotal').innerText = s.total;
      document.getElementById('metricTasksRun').innerText = s.running;
      document.getElementById('metricTasksSleep').innerText = s.sleeping;
      document.getElementById('metricTasksStop').innerText = s.stopped;
      document.getElementById('metricTasksZombie').innerText = s.zombie;

      document.getElementById('cpuUserVal').innerText = (m.cpu.user || 0).toFixed(1);
      document.getElementById('cpuSysVal').innerText = (m.cpu.system || 0).toFixed(1);
      document.getElementById('cpuIdleVal').innerText = (m.cpu.idle || 0).toFixed(1);
      document.getElementById('cpuWaitVal').innerText = (m.cpu.iowait || 0).toFixed(1);
      document.getElementById('cpuTotalText').innerText = (m.cpu.overall || 0).toFixed(1) + '%%';
      document.getElementById('cpuTotalBar').style.width = Math.min(100, m.cpu.overall || 0) + '%%';

      const mem = m.memory;
      document.getElementById('memTotalTxt').innerText = formatByteUnits(mem.total);
      document.getElementById('memFreeTxt').innerText = formatByteUnits(mem.free);
      document.getElementById('memUsedTxt').innerText = formatByteUnits(mem.used);
      document.getElementById('memCacheTxt').innerText = formatByteUnits((mem.cached || 0) + (mem.buffers || 0));
      document.getElementById('memUsedBar').style.width = (mem.used_percent || 0) + '%%';
      document.getElementById('memCacheBar').style.width = Math.min(100, (((mem.cached || 0) / mem.total) * 100)) + '%%';

      document.getElementById('swapTotalTxt').innerText = formatByteUnits(mem.swap_total);
      document.getElementById('swapFreeTxt').innerText = formatByteUnits(mem.swap_free);
      document.getElementById('swapUsedTxt').innerText = formatByteUnits(mem.swap_used);
      document.getElementById('swapUsedBar').style.width = (mem.swap_percent || 0) + '%%';

      sparklineHistory.push(m.cpu.overall || 0);
      sparklineHistory.shift();
      const sparkGrid = document.getElementById('sparklineGrid');
      sparkGrid.innerHTML = sparklineHistory.map(v => 
        '<div class="spark-segment" style="height:' + Math.max(4, Math.min(100, v)) + '%%;' + (v > 80 ? 'background:var(--color-danger)' : v > 55 ? 'background:var(--color-warning)' : '') + '"></div>'
      ).join('');

      const cores = m.cpu.per_core || [m.cpu.overall];
      const coresContainer = document.getElementById('coresContainer');
      coresContainer.innerHTML = cores.map((c, i) => 
        '<div class="core-cell">' +
          '<div style="display:flex; justify-content:space-between; font-size:11px; margin-bottom:4px;">' +
            '<span>CPU ' + i + '</span><strong style="color:var(--color-accent-green)">' + c.toFixed(1) + '%%</strong>' +
          '</div>' +
          '<div class="bar-track"><div class="bar-fill" style="width:' + Math.min(100, c) + '%%; background:' + (c > 80 ? 'var(--color-danger)' : c > 55 ? 'var(--color-warning)' : 'var(--color-accent-green)') + '"></div></div>' +
        '</div>'
      ).join('');

      const disks = m.disks || [];
      document.getElementById('mountsContainer').innerHTML = disks.map(d => 
        '<div style="margin-bottom:12px;">' +
          '<div style="display:flex; justify-content:space-between; font-size:11px;">' +
            '<span><strong>' + d.mount_point + '</strong> (' + d.device + ')</span>' +
            '<span style="color:var(--color-cyan)">' + d.used_percent + '%%</span>' +
          '</div>' +
          '<div class="bar-track"><div class="bar-fill" style="background:var(--color-cyan); width:' + d.used_percent + '%%"></div></div>' +
          '<div style="display:flex; justify-content:space-between; font-size:10px; color:var(--text-secondary); margin-top:3px;">' +
            '<span>' + formatByteUnits(d.used_bytes) + ' used</span>' +
            '<span>' + formatByteUnits(d.free_bytes) + ' free</span>' +
          '</div>' +
        '</div>'
      ).join('');

      const netList = Object.values(m.network || {});
      document.getElementById('netContainer').innerHTML = netList.map(n => 
        '<div style="padding:8px 0; border-bottom:1px solid rgba(255,255,255,0.04); font-size:11px;">' +
          '<div style="display:flex; justify-content:space-between; font-weight:600;">' +
            '<span>' + n.name + '</span>' +
            '<span style="color:var(--text-secondary)">' + (n.rx_packets + n.tx_packets) + ' pkts</span>' +
          '</div>' +
          '<div style="display:flex; justify-content:space-between; margin-top:3px; font-size:10px;">' +
            '<span style="color:var(--color-accent-green)">⬇ RX: ' + formatByteUnits(n.rx_rate_bytes) + '/s</span>' +
            '<span style="color:var(--color-cyan)">⬆ TX: ' + formatByteUnits(n.tx_rate_bytes) + '/s</span>' +
          '</div>' +
        '</div>'
      ).join('');
    }

    function setFilterState(state) {
      stateFilterMode = state;
      ['ALL', 'R', 'S', 'Z'].forEach(k => {
        document.getElementById('btn-' + k).className = (k === state) ? 'primary' : '';
      });
      applyFilters();
    }

    function toggleSort(col) {
      if (sortColumnKey === col) {
        isSortAscending = !isSortAscending;
      } else {
        sortColumnKey = col;
        isSortAscending = false;
      }
      applyFilters();
    }

    function applyFilters() {
      const q = (document.getElementById('searchInput').value || '').toLowerCase();
      let list = processRecords.filter(p => {
        if (stateFilterMode !== 'ALL' && p.state !== stateFilterMode) return false;
        if (!q) return true;
        return p.name.toLowerCase().includes(q) ||
               (p.cmdline || '').toLowerCase().includes(q) ||
               p.pid.toString().includes(q) ||
               p.user.toLowerCase().includes(q);
      });

      list.sort((a, b) => {
        let vA = a[sortColumnKey], vB = b[sortColumnKey];
        if (typeof vA === 'string') return isSortAscending ? vA.localeCompare(vB) : vB.localeCompare(vA);
        return isSortAscending ? (vA - vB) : (vB - vA);
      });

      const tbody = document.getElementById('processesTbody');
      if (list.length === 0) {
        tbody.innerHTML = '<tr><td colspan="12" style="text-align:center; color:var(--text-secondary); padding:24px;">No matching processes found.</td></tr>';
        return;
      }

      tbody.innerHTML = list.slice(0, 100).map(p => 
        '<tr>' +
          '<td><strong>' + p.pid + '</strong></td>' +
          '<td>' + p.user + '</td>' +
          '<td>' + p.priority + '</td>' +
          '<td>' + p.nice + '</td>' +
          '<td>' + formatByteUnits(p.virt_bytes) + '</td>' +
          '<td>' + formatByteUnits(p.res_bytes) + '</td>' +
          '<td><span class="proc-state state-' + p.state + '">' + p.state + '</span></td>' +
          '<td style="color:' + (p.cpu_percent > 30 ? 'var(--color-danger)' : 'var(--color-accent-green)') + '; font-weight:700;">' + p.cpu_percent.toFixed(1) + '</td>' +
          '<td>' + p.mem_percent.toFixed(1) + '</td>' +
          '<td>' + p.time_string + '</td>' +
          '<td style="max-width:320px; overflow:hidden; text-overflow:ellipsis; white-space:nowrap;" title="' + (p.cmdline || p.name) + '">' +
            '<strong>' + p.name + '</strong> <span style="color:var(--text-muted); font-size:11px;">' + (p.cmdline !== p.name ? p.cmdline : '') + '</span>' +
          '</td>' +
          '<td>' +
            '<button class="danger" style="padding:2px 6px; font-size:10px;" onclick="terminateProcess(' + p.pid + ')">SIGTERM</button>' +
          '</td>' +
        '</tr>'
      ).join('');
    }

    async function terminateProcess(pid) {
      if (!confirm('Send SIGTERM to process PID ' + pid + '?')) return;
      try {
        const res = await fetch('/api/process/signal', {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          body: JSON.stringify({ pid, signal: 'SIGTERM' })
        });
        const data = await res.json();
        alert(data.message || 'Signal processed.');
        pollTelemetry();
      } catch(e) {
        alert('Signal error: ' + e);
      }
    }

    async function execSystemd(action) {
      try {
        const res = await fetch('/api/systemd/action', {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          body: JSON.stringify({ action })
        });
        const data = await res.json();
        alert('Action ' + action + ' executed.');
        fetchJournalLogs();
      } catch(e) {
        alert('Systemd error: ' + e);
      }
    }

    async function fetchJournalLogs() {
      try {
        const res = await fetch('/api/systemd/logs');
        const data = await res.json();
        const box = document.getElementById('journalOutput');
        box.innerHTML = (data.logs || []).map(l => 
          '<div class="log-row"><span style="color:var(--text-secondary)">[' + new Date(l.timestamp).toLocaleTimeString() + ']</span> <strong style="color:var(--color-accent-green)">' + l.level + '</strong> ' + l.message + '</div>'
        ).join('');
      } catch(e){}
    }

    async function fetchRaw(path) {
      try {
        const res = await fetch('/api/proc-raw?file=' + encodeURIComponent(path));
        const txt = await res.text();
        document.getElementById('rawBox').innerText = txt;
      } catch(e) {
        document.getElementById('rawBox').innerText = 'Failed to load ' + path;
      }
    }

    setInterval(pollTelemetry, 1000);
    pollTelemetry();
  </script>
</body>
</html>`
