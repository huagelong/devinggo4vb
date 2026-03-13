import React from 'react'
import ReactDOM from 'react-dom/client'
import { RouterProvider } from '@tanstack/react-router'
import { ConfigProvider, theme } from 'antd'
import zhCN from 'antd/locale/zh_CN'
import { router } from './router'
import './index.css'
import './i18n'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <ConfigProvider
      locale={zhCN}
      theme={{
        algorithm: theme.compactAlgorithm,
        token: {
          colorPrimary: '#1677ff',
          colorBgContainer: '#ffffff',
          borderRadius: 4,
          fontSize: 13,
          wireframe: true,
          fontFamily: 'Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif'
        },
        components: {
          Button: {
            controlHeight: 28,
            controlHeightSM: 24,
            controlHeightLG: 32,
            paddingInline: 12,
          },
          Input: {
            controlHeight: 28,
            controlHeightLG: 32,
          },
          Select: {
            controlHeight: 28,
            controlHeightLG: 32,
          },
          Table: {
            padding: 12,
            paddingSM: 8,
          },
        },
      }}
    >
      <RouterProvider router={router} />
    </ConfigProvider>
  </React.StrictMode>,
)
