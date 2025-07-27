package main

// Language represents the application language
type Language string

const (
	LangEnglish Language = "en"
	LangTurkish Language = "tr"
)

// Localization holds all text translations
type Localization struct {
	// App titles
	NamespaceSelectionTitle string
	LogAnalysisTitle        string
	NamespaceTitle          string
	Pods                    string

	// Pod states
	PodDetails   string
	Name         string
	Status       string
	Ready        string
	Restart      string
	Age          string
	Analysis     string
	LogSummary   string
	TotalLines   string
	Errors       string
	Warnings     string
	LogLines     string
	ShowingLines string
	TotalFrom    string
	LastLines    string

	// Status messages
	NamespaceNotFound string
	PodNotFound       string
	LogNotFound       string
	Loading           string
	LogEmpty          string
	StatusNormal      string
	StatusWarning     string
	StatusError       string

	// Navigation
	Controls          string
	Movement          string
	Select            string
	ViewLogs          string
	GoBack            string
	Refresh           string
	AutoRefresh       string
	AutoRefreshStatus string
	Exit              string
	RefreshLogs       string
	UpDown            string
	LeftRight         string
	ScrollUp          string
	ScrollDown        string

	// Pagination
	Namespaces string
}

// GetLocalization returns localization based on language
func GetLocalization(lang Language) Localization {
	switch lang {
	case LangTurkish:
		return Localization{
			// App titles
			NamespaceSelectionTitle: "Kubernetes Namespace Seçimi",
			LogAnalysisTitle:        "Log Analizi",
			NamespaceTitle:          "Namespace",
			Pods:                    "Pod'lar",

			// Pod states
			PodDetails:   "Pod Detayları",
			Name:         "İsim",
			Status:       "Durum",
			Ready:        "Hazır",
			Restart:      "Restart",
			Age:          "Yaş",
			Analysis:     "Analiz",
			LogSummary:   "Log Özeti",
			TotalLines:   "Toplam satır",
			Errors:       "Hatalar",
			Warnings:     "Uyarılar",
			LogLines:     "Log Satırları",
			ShowingLines: "satır gösteriliyor",
			TotalFrom:    "Toplam",
			LastLines:    "satırdan son",

			// Status messages
			NamespaceNotFound: "Namespace bulunamadı",
			PodNotFound:       "Pod bulunamadı",
			LogNotFound:       "Log analizi bulunamadı",
			Loading:           "Yükleniyor...",
			LogEmpty:          "Log bulunamadı veya boş",
			StatusNormal:      "DURUM: Normal",
			StatusWarning:     "DURUM: Uyarı var",
			StatusError:       "DURUM: Hata var",

			// Navigation
			Controls:          "Kontroller",
			Movement:          "k/j veya ok tuşları: Hareket",
			Select:            "Enter: Namespace seç",
			ViewLogs:          "Enter: Log görüntüle",
			GoBack:            "Esc/Backspace: Geri dön",
			Refresh:           "r: Yenile",
			AutoRefresh:       "t: Auto-refresh",
			AutoRefreshStatus: "Otomatik yenileme",
			Exit:              "q: Çıkış",
			RefreshLogs:       "r: Logları yenile",
			UpDown:            "Yukarı/Aşağı: k/j veya ok tuşları",
			LeftRight:         "Sol/Sağ: h/l veya ok tuşları",
			ScrollUp:          "Yukarı kaydır",
			ScrollDown:        "Aşağı kaydır",

			// Pagination
			Namespaces: "namespace",
		}
	default: // English
		return Localization{
			// App titles
			NamespaceSelectionTitle: "Kubernetes Namespace Selection",
			LogAnalysisTitle:        "Log Analysis",
			NamespaceTitle:          "Namespace",
			Pods:                    "Pods",

			// Pod states
			PodDetails:   "Pod Details",
			Name:         "Name",
			Status:       "Status",
			Ready:        "Ready",
			Restart:      "Restart",
			Age:          "Age",
			Analysis:     "Analysis",
			LogSummary:   "Log Summary",
			TotalLines:   "Total lines",
			Errors:       "Errors",
			Warnings:     "Warnings",
			LogLines:     "Log Lines",
			ShowingLines: "lines showing",
			TotalFrom:    "Total",
			LastLines:    "last lines from",

			// Status messages
			NamespaceNotFound: "Namespace not found",
			PodNotFound:       "Pod not found",
			LogNotFound:       "Log analysis not found",
			Loading:           "Loading...",
			LogEmpty:          "Log not found or empty",
			StatusNormal:      "STATUS: Normal",
			StatusWarning:     "STATUS: Warning",
			StatusError:       "STATUS: Error",

			// Navigation
			Controls:          "Controls",
			Movement:          "k/j or arrow keys: Move",
			Select:            "Enter: Select namespace",
			ViewLogs:          "Enter: View logs",
			GoBack:            "Esc/Backspace: Go back",
			Refresh:           "r: Refresh",
			AutoRefresh:       "t: Auto-refresh",
			AutoRefreshStatus: "Auto-refresh",
			Exit:              "q: Exit",
			RefreshLogs:       "r: Refresh logs",
			UpDown:            "Up/Down: k/j or arrow keys",
			LeftRight:         "Left/Right: h/l or arrow keys",
			ScrollUp:          "Scroll up",
			ScrollDown:        "Scroll down",

			// Pagination
			Namespaces: "namespaces",
		}
	}
}
