// Package i18n
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package i18n

import (
	"context"
	"devinggo/modules/system/pkg/contexts"
	"devinggo/modules/system/pkg/utils/request"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
)

const (
	DefaultLanguage = "zh-CN"
)

var (
	once sync.Once
)

// initI18nPath 初始化国际化资源路径（只需要初始化一次）
func initI18nPath() {
	once.Do(func() {
		g.I18n().SetPath("resource/i18n")
	})
}

// GetLanguageFromRequest 从请求中获取语言设置
func GetLanguageFromRequest(ctx context.Context) string {
	r := request.GetHttpRequest(ctx)

	// 优先级1: URL参数 lang
	langGet := r.Get("lang")
	if !langGet.IsEmpty() {
		return normalizeLanguage(langGet.String())
	}

	// 优先级2: Header Accept-Language
	headerLang := r.Header.Get("Accept-Language")
	if !g.IsEmpty(headerLang) {
		return normalizeLanguage(headerLang)
	}

	// 默认语言
	return DefaultLanguage
}

// normalizeLanguage 规范化语言标识
func normalizeLanguage(lang string) string {
	// 处理 "zh-CN;q=0.9,en;q=0.8" 这样的格式
	if gstr.Contains(lang, ";") {
		lang = gstr.Split(lang, ";")[0]
	}

	// 处理 "zh-CN,zh;q=0.9" 这样的格式
	if gstr.Contains(lang, ",") {
		lang = gstr.Split(lang, ",")[0]
	}

	lang = gstr.Trim(lang)

	// 简化英文语言标识
	if gstr.ContainsI(lang, "en") {
		return "en"
	}

	return lang
}

// InitI18n 初始化国际化配置
func InitI18n(ctx context.Context) {
	initI18nPath()

	lang := GetLanguageFromRequest(ctx)
	g.I18n().SetLanguage(lang)

	// 将语言信息存入上下文，方便后续使用
	contexts.SetLanguage(ctx, lang)
}

// T 翻译文本（简单版本）
func T(ctx context.Context, key string) string {
	initI18nPath()
	return g.I18n().T(ctx, key)
}

// Tf 翻译文本并格式化（接受任意类型参数）
func Tf(ctx context.Context, key string, params ...interface{}) string {
	initI18nPath()
	return g.I18n().Tf(ctx, key, params...)
}

// GetCurrentLanguage 获取当前语言
func GetCurrentLanguage(ctx context.Context) string {
	if lang := contexts.GetLanguage(ctx); lang != "" {
		return lang
	}
	return DefaultLanguage
}
