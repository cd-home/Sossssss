package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(theme.LightTheme())
	//appIcon, _ := fyne.LoadResourceFromPath("./icon/app.png")
	//myApp.SetIcon(appIcon)
	myWindow := myApp.NewWindow("工具")
	// 窗口
	myWindow.Resize(fyne.NewSize(1000, 540))
	// 居中
	myWindow.CenterOnScreen()
	// tab
	tab1 := widget.NewLabel("Hello")
	tab2 := widget.NewCard("title", "subtitle", widget.NewLabel("123"))
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Home", theme.FolderOpenIcon(), widget.NewLabel("Home tab")),
		container.NewTabItem("Tab 1", tab1),
		container.NewTabItem("Tab 2", tab2),
	)
	tabs.SetTabLocation(container.TabLocationLeading)
	tabs.DisableIndex(0)

	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}

type myTheme struct {
	d theme.ThemedResource
}

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {
		if variant == theme.VariantLight {
			return color.White
		}
		return color.Black
	}
	return theme.DefaultTheme().Color(name, variant)
}
