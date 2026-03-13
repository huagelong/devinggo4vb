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
          fontSize: 17,
          wireframe: true,
          fontFamily: 'Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif'
        },
        components: {
          Button: {
            controlHeight: 38,
            controlHeightSM: 34,
            controlHeightLG: 43,
            paddingInline: 17,
          },
          Input: {
            controlHeight: 38,
            controlHeightLG: 43,
          },
          Select: {
            controlHeight: 38,
            controlHeightLG: 43,
          },
          Table: {
            padding: 14,
            paddingSM: 10,
          },
        },
      }}
    >
      <RouterProvider router={router} />
    </ConfigProvider>
  </React.StrictMode>,
)
