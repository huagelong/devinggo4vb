import React, { useMemo } from 'react';
import { ProLayout, MenuDataItem } from '@ant-design/pro-components';
import { Outlet, useNavigate, useLocation } from '@tanstack/react-router';
import { Dropdown, MenuProps, Space } from 'antd';
import { LogoutOutlined, UserOutlined, SettingOutlined, DashboardOutlined, SettingFilled, AppstoreOutlined } from '@ant-design/icons';
import { useAuthStore } from '../../stores/authStore';
import { useTranslation } from 'react-i18next';

export const BasicLayout: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { userInfo, routers, logout } = useAuthStore();
  const { t } = useTranslation();

  const handleLogout = () => {
    logout();
    navigate({ to: '/login' });
  };

  const userMenuItems: MenuProps['items'] = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人中心',
    },
    {
      key: 'settings',
      icon: <SettingOutlined />,
      label: '个人设置',
    },
    {
      type: 'divider',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: t('system.logout'),
      onClick: handleLogout,
      danger: true,
    },
  ];

  // 动态将后端菜单数据转换为 ProLayout 所需的 Route 结构
  const menuData = useMemo(() => {
    const mapRoutersToMenuData = (apiRouters: any[], parentPath = ''): MenuDataItem[] => {
      if (!apiRouters) return [];
      return apiRouters.map(r => {
        let currentPath = parentPath ? `${parentPath}${r.path.startsWith('/') ? r.path : '/' + r.path}` : r.path;
        currentPath = currentPath.replace(/\/+/g, '/'); 
        
        // 适配本地临时路由结构: 将接口返回的 systemMg 统一映射到本地的 system 主干路由
        currentPath = currentPath.replace(/^\/systemMg/, '/system');

        return {
          path: currentPath,
          name: r.meta?.title || r.name,
          // 根据 type=='B' 判断是否为按钮级别权限，不需要渲染为侧边栏菜单
          hideInMenu: r.meta?.hidden === true || r.meta?.type === 'B',
          routes: r.children && r.children.length > 0 ? mapRoutersToMenuData(r.children, currentPath) : undefined,
        };
      });
    };

    const dynamicMenus = routers ? mapRoutersToMenuData(routers) : [];

    // 添加固定的看板导航
    return {
      path: '/',
      routes: [
        {
          path: '/dashboard',
          name: '工作台',
          icon: <DashboardOutlined />
        },
        ...dynamicMenus,
      ],
    };
  }, [routers]);

  return (
    <div style={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
      <ProLayout
        title="Admin Pro"
        logo={<img src="https://gw.alipayobjects.com/zos/rmsportal/KDpgvguMpGfqaHPjicRK.svg" alt="logo" />}
        layout="mix"
        splitMenus={false}
        fixSiderbar
        fixedHeader
        siderMenuType="group"
        token={{
          bgLayout: '#f0f2f5',
          header: {
            colorBgHeader: '#fff',
            colorHeaderTitle: '#141414',
            colorTextMenu: '#dfdfdf',
            colorTextMenuSecondary: '#dfdfdf',
            colorTextMenuSelected: '#fff',
            colorBgMenuItemSelected: '#22272b',
            colorTextMenuActive: '#rgba(255,255,255,0.85)',
          },
          sider: {
            colorMenuBackground: '#fff',
            colorMenuItemDivider: '#dfdfdf',
            colorTextMenu: '#595959',
            colorTextMenuSelected: '#1677ff',
            colorBgMenuItemSelected: '#e6f4ff',
          },
          pageContainer: {
            colorBgPageContainer: '#f0f2f5',
            paddingBlockPageContainerContent: 24,
            paddingInlinePageContainerContent: 24,
          }
        }}
        avatarProps={{
          src: userInfo?.avatar || 'https://gw.alipayobjects.com/zos/antfincdn/efFD%24IOql2/weixintupian_20170331104822.jpg',
          title: userInfo?.username || 'Admin',
          size: 'small',
          render: (_props, dom) => {
            return (
              <Dropdown menu={{ items: userMenuItems }} trigger={['click']}>
                <div className="flex items-center cursor-pointer hover:bg-gray-100 px-3 rounded transition-colors h-full">
                  {dom}
                </div>
              </Dropdown>
            );
          },
        }}
        actionsRender={() => [
          <div key="actions" className="flex items-center pr-4">
            <span className="text-gray-500 hover:text-gray-900 cursor-pointer text-sm">EN</span>
          </div>
        ]}
        route={menuData}
        location={{
          pathname: location.pathname,
        }}
        menuItemRender={(item, dom) => (
          <a
            onClick={(e) => {
              e.preventDefault();
              if (item.path) {
                navigate({ to: item.path });
              }
            }}
          >
            {dom}
          </a>
        )}
      >
        <div style={{ minHeight: 'calc(100vh - 110px)', padding: '24px' }}>
          <Outlet />
        </div>
      </ProLayout>
    </div>
  );
};

export default BasicLayout;
