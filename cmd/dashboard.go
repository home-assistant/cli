package cmd

import (
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/thebanri/limoni"
	"github.com/thebanri/limoni/widgets"
)

type dashboardMetrics struct {
	CPUPercent    float64
	MemoryPercent float64
	DiskPercent   float64
}

type dashboardAddon struct {
	Name    string
	Running bool
}

type dashboardState struct {
	mu sync.RWMutex

	Connected bool
	LastError string
	LastFetch time.Time

	Metrics dashboardMetrics
	Addons  []dashboardAddon
	Logs    []string
}

func (s *dashboardState) snapshot() dashboardState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return dashboardState{
		Connected: s.Connected,
		LastError: s.LastError,
		LastFetch: s.LastFetch,
		Metrics:   s.Metrics,
		Addons:    slices.Clone(s.Addons),
		Logs:      slices.Clone(s.Logs),
	}
}

func (s *dashboardState) setFromFetch(connected bool, err error, m dashboardMetrics, addons []dashboardAddon, logs []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Connected = connected
	if err != nil {
		s.LastError = err.Error()
	} else {
		s.LastError = ""
	}
	s.LastFetch = time.Now()
	s.Metrics = m
	s.Addons = slices.Clone(addons)

	const maxLogs = 200
	if len(logs) > maxLogs {
		logs = logs[len(logs)-maxLogs:]
	}
	s.Logs = slices.Clone(logs)
}

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Launch interactive Home Assistant TUI dashboard",
	Long:  "Launches an interactive, real-time terminal dashboard for Home Assistant Supervisor telemetry and status.",
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		state := &dashboardState{}
		refreshReq := make(chan struct{}, 1)
		stopFetch := make(chan struct{})
		defer close(stopFetch)

		// Trigger immediate initial fetch
		refreshReq <- struct{}{}

		go func() {
			ticker := time.NewTicker(1 * time.Second) // Saatin saniye saniye akması için 1sn
			defer ticker.Stop()

			for {
				select {
					case <-stopFetch:
						return
					case <-ticker.C:
						fetchAndStoreDashboardData(state)
						limoni.Wakeup() // <-- BURAYA: Render döngüsünü uyandırır
					case <-refreshReq:
						fetchAndStoreDashboardData(state)
						limoni.Wakeup() // <-- BURAYA: Anında ekranı yeniler
				}
			}
		}()


		return limoni.Run(func(f *limoni.Frame, ev *limoni.Event) bool {
			if ev != nil && ev.Type == limoni.EventKey {
				switch ev.Key.Type {
				case limoni.KeyEsc:
					return false
				case limoni.KeyRune:
					switch ev.Key.Ch {
					case 'q', 'Q':
						return false
					case 'r', 'R':
						select {
						case refreshReq <- struct{}{}:
						default:
						}
					}
				}
			}

			s := state.snapshot()
			renderDashboard(f, s)
			return true
		})
	},
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
}

func fetchAndStoreDashboardData(state *dashboardState) {
	metrics, addons, logs, err := fetchDashboardData()
	if err != nil {
		// Mock veriyi mevcut log durumunu koruyarak al
		state.mu.RLock()
		currentLogs := state.Logs
		state.mu.RUnlock()

		metrics, addons, logs = mockDashboardData(err, currentLogs)
		state.setFromFetch(false, err, metrics, addons, logs)
		return
	}
	state.setFromFetch(true, nil, metrics, addons, logs)
}

func fetchDashboardData() (dashboardMetrics, []dashboardAddon, []string, error) {
	// TODO: Replace with existing client/supervisor calls in ha CLI
	return dashboardMetrics{}, nil, nil, fmt.Errorf("supervisor data source not yet wired")
}


var initialMockLogs = []string{
	"2026-09-10T13:20:01Z [info] supervisor: development fallback mode enabled",
	"2026-09-10T13:20:02Z [warn] api: unable to reach supervisor, using mocked telemetry",
	"2026-09-10T13:20:05Z [info] core: Home Assistant healthy",
	"2026-09-10T13:20:10Z [info] jobs: scheduler idle",
}

func mockDashboardData(fetchErr error, existingLogs []string) (dashboardMetrics, []dashboardAddon, []string) {
	metrics := dashboardMetrics{
		CPUPercent:    23.4,
		MemoryPercent: 61.2,
		DiskPercent:   47.8,
	}
	addons := []dashboardAddon{
		{Name: "Terminal & SSH", Running: true},
		{Name: "File editor", Running: true},
		{Name: "Mosquitto broker", Running: true},
		{Name: "Samba share", Running: false},
		{Name: "Node-RED", Running: false},
	}

	// İlk açılışta hazır logları koy
	if len(existingLogs) == 0 {
		existingLogs = slices.Clone(initialMockLogs)
	}

	return metrics, addons, existingLogs
}

func renderDashboard(f *limoni.Frame, s dashboardState) {
	now := time.Now().Format("2006-01-02 15:04:05")

	statusLabel := limoni.Label("[DISCONNECTED]", limoni.Fg(limoni.Hex("#FF5555")).Bold())
	if s.Connected {
		statusLabel = limoni.Label("[CONNECTED]", limoni.Fg(limoni.Hex("#00FFAA")).Bold())
	}

	header := limoni.FixedSize(0, 3, limoni.Border(
		limoni.Pad(
			limoni.HStack(
				limoni.Label("HOME ASSISTANT MONITOR", limoni.Bold().WithFg(limoni.Hex("#00FFAA"))),
				statusLabel,
				limoni.Label(now, limoni.Fg(limoni.Hex("#88CCFF"))),
			).WithJustify(limoni.JustifySpaceBetween).WithAlignItems(limoni.AlignItemsCenter),
			0, 1, 0, 1,
		),
		widgets.SymbolsRounded,
		limoni.Fg(limoni.Hex("#00FFAA")),
	))

	cpuBox := limoni.Flex(1, limoni.Border(
		limoni.Pad(
			limoni.VStack(
				limoni.Label("CPU Usage", limoni.Fg(limoni.Hex("#AAAAAA"))),
				limoni.Label(fmt.Sprintf("%.1f%%", s.Metrics.CPUPercent), limoni.Bold().WithFg(limoni.Hex("#00FFAA"))),
			),
			0, 1, 0, 1,
		),
		widgets.SymbolsRounded,
		limoni.Fg(limoni.Hex("#00FFAA")),
	))

	memBox := limoni.Flex(1, limoni.Border(
		limoni.Pad(
			limoni.VStack(
				limoni.Label("Memory", limoni.Fg(limoni.Hex("#AAAAAA"))),
				limoni.Label(fmt.Sprintf("%.1f%%", s.Metrics.MemoryPercent), limoni.Bold().WithFg(limoni.Hex("#3399FF"))),
			),
			0, 1, 0, 1,
		),
		widgets.SymbolsRounded,
		limoni.Fg(limoni.Hex("#3399FF")),
	))

	diskBox := limoni.Flex(1, limoni.Border(
		limoni.Pad(
			limoni.VStack(
				limoni.Label("Disk Storage", limoni.Fg(limoni.Hex("#AAAAAA"))),
				limoni.Label(fmt.Sprintf("%.1f%%", s.Metrics.DiskPercent), limoni.Bold().WithFg(limoni.Hex("#FFCC00"))),
			),
			0, 1, 0, 1,
		),
		widgets.SymbolsRounded,
		limoni.Fg(limoni.Hex("#FFCC00")),
	))

	metricsPanel := limoni.FixedSize(0, 4, limoni.HStack(cpuBox, memBox, diskBox).WithGap(1))

	addonItems := make([]limoni.Component, 0, len(s.Addons)+1)
	addonItems = append(addonItems, limoni.Label("Managed Add-ons", limoni.Bold().WithFg(limoni.Hex("#FFFFFF"))))
	for _, a := range s.Addons {
		statusColor := limoni.Hex("#FF5555")
		statusText := "STOPPED"
		if a.Running {
			statusColor = limoni.Hex("#00FFAA")
			statusText = "RUNNING"
		}
		row := limoni.HStack(
			limoni.Label(fmt.Sprintf("- %s", a.Name), limoni.Fg(limoni.Hex("#CCCCCC"))),
			limoni.Label(fmt.Sprintf("[%s]", statusText), limoni.Bold().WithFg(statusColor)),
		).WithJustify(limoni.JustifySpaceBetween)
		addonItems = append(addonItems, row)
	}

	addonsPane := limoni.Border(
		limoni.Pad(limoni.VStack(addonItems...), 1, 1, 1, 1),
		widgets.SymbolsRounded,
		limoni.Fg(limoni.Hex("#5588EE")),
	)

	logLines := s.Logs
	if len(logLines) == 0 {
		logLines = []string{"No log events yet."}
	}
	logItems := make([]limoni.Component, 0, len(logLines)+1)
	logItems = append(logItems, limoni.Label("Supervisor Logs & Telemetry", limoni.Bold().WithFg(limoni.Hex("#FFFFFF"))))
	for _, line := range logLines {
		logItems = append(logItems, limoni.Label(line, limoni.Fg(limoni.Hex("#888888"))))
	}

	logPane := limoni.Border(
		limoni.Pad(limoni.VStack(logItems...), 1, 1, 1, 1),
		widgets.SymbolsRounded,
		limoni.Fg(limoni.Hex("#3399FF")),
	)

	body := limoni.Flex(1, limoni.HStack(
		limoni.Flex(1, addonsPane),
		limoni.Flex(2, logPane),
	).WithGap(1))

	footer := limoni.FixedSize(0, 3, limoni.Border(
		limoni.Pad(
			limoni.HStack(
				limoni.Label("r: Refresh | q / ESC: Exit", limoni.Fg(limoni.Hex("#888888"))),
				limoni.Label("Pure Go | Limoni 0 B/op", limoni.Bold().WithFg(limoni.Hex("#00FFAA"))),
			).WithJustify(limoni.JustifySpaceBetween).WithAlignItems(limoni.AlignItemsCenter),
			0, 1, 0, 1,
		),
		widgets.SymbolsSingle,
		limoni.Fg(limoni.Hex("#555555")),
	))

	rootView := limoni.VStack(
		header,
		metricsPanel,
		body,
		footer,
	).WithGap(1)

	f.RenderComponent(rootView, f.Area())
}
