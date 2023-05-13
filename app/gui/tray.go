package main

import (
	"github.com/getlantern/systray"
	"stroxy/boot"
	"stroxy/config"
)

// 启动GUI系统托盘
func startUp() {
	systray.Run(onReady, onExit)
}

// 在切换节点选项下为每个节点添加按钮并设置点击事件
func addNodeItem(mNode *systray.MenuItem) {
	nodeMap := config.ProductConfigGroup.ProxyServerMap
	var defaultNodeMenu *systray.MenuItem
	for k := range nodeMap {
		curNodeMenu := mNode.AddSubMenuItem(k, k)

		if config.ProductConfigGroup.DefaultServer == k {
			curNodeMenu.Disable()
			defaultNodeMenu = curNodeMenu
		}

		go func(name string, menu *systray.MenuItem) {
			for range curNodeMenu.ClickedCh {
				boot.SelectNode(name)
				defaultNodeMenu.Enable()
				menu.Disable()
				defaultNodeMenu = menu
			}
		}(k, curNodeMenu)
	}

}

// 添加代理模式按钮，auto\all
func addProxyModeItem() (autoProxyMenu, allProxyMenu *systray.MenuItem) {
	mProxyMode := systray.AddMenuItem("代理模式", "代理模式")
	allProxyMenu = mProxyMode.AddSubMenuItem("全局代理", "全局代理")
	autoProxyMenu = mProxyMode.AddSubMenuItem("智能代理", "智能代理")
	switch config.ProductConfigGroup.ProxyModel {
	case "auto":
		autoProxyMenu.Disable()
	case "all":
		allProxyMenu.Disable()
	}
	return
}

func onReady() {
	systray.SetTemplateIcon(Data, Data)

	systray.SetTitle("Stroxy")
	systray.SetTooltip("Stroxy")

	mRun := systray.AddMenuItem("启动", "启动")

	mStop := systray.AddMenuItem("暂停", "暂停")
	mNode := systray.AddMenuItem("切换节点", "切换节点")
	addNodeItem(mNode)

	autoProxyMenu, allProxyMenu := addProxyModeItem()

	mQuit := systray.AddMenuItem("退出", "退出程序")

	go func() {
		for {
			select {
			case <-mRun.ClickedCh:
				boot.RunProxy()
				mRun.Disable()
				mNode.Disable()
				mStop.Enable()
			case <-mStop.ClickedCh:
				boot.StopProxy()
				mRun.Enable()
				mNode.Enable()
				mStop.Disable()
			case <-autoProxyMenu.ClickedCh:
				autoProxyMenu.Disable()
				allProxyMenu.Enable()
				boot.SelectProxyMode("auto")
			case <-allProxyMenu.ClickedCh:
				allProxyMenu.Disable()
				autoProxyMenu.Enable()
				boot.SelectProxyMode("all")
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
	mRun.ClickedCh <- struct{}{}
}

func onExit() {
	boot.ExitSystem()
}
