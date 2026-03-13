import React, { useMemo, useState } from 'react';
import { ProLayout, MenuDataItem } from '@ant-design/pro-components';
import { Outlet, useNavigate, useLocation } from '@tanstack/react-router';
import { Dropdown, MenuProps, Avatar, Badge } from 'antd';
import { 
  LogoutOutlined, UserOutlined, SettingOutlined, DashboardOutlined, 
  HomeOutlined, AppstoreOutlined, SearchOutlined, LockOutlined, 
  FullscreenOutlined, BellOutlined, CloseOutlined, MenuUnfoldOutlined, MenuFoldOutlined
} from '@ant-design/icons';
import { useAuthStore } from '../../stores/authStore';
import { useTranslation } from 'react-i18next';

export const BasicLayout: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { userInfo, routers, logout } = useAuthStore();
  const { t } = useTranslation();
  
  const [collapsed, setCollapsed] = useState(false);

  const handleLogout = () => {
    logout();
    navigate({ to: '/login' });
  };

  const userMenuItems: MenuProps['items'] = [
    { key: 'profile', icon: <UserOutlined />, label: '个人中心' },
    { key: 'settings', icon: <SettingOutlined />, label: '个人设置' },
    { type: 'divider' },
    { key: 'logout', icon: <LogoutOutlined />, label: t('system.logout'), onClick: handleLogout, danger: true },
  ];

  const fullMenuData = useMemo(() => {
    const mapRoutersToMenuData = (apiRouters: any[], parentPath = ''): MenuDataItem[] => {
      if (!apiRouters) return [];
      return apiRouters.map(r => {
        let currentPath = parentPath ? `${parentPath}${r.path.startsWith('/') ? r.path : '/' + r.path}` : r.path;
        currentPath = currentPath.replace(/\/+/g, '/');
        currentPath = currentPath.replace(/^\/systemMg/, '/system');

        return {
          path: currentPath,
          name: r.meta?.title || r.name,
          hideInMenu: r.meta?.hidden === true || r.meta?.type === 'B',
          routes: r.children && r.children.length > 0 ? mapRoutersToMenuData(r.children, currentPath) : undefined,
          icon: <AppstoreOutlined />,
        };
      });
    };

    const dynamicMenus = routers ? mapRoutersToMenuData(routers) : [];

    return [
      {
        path: '/dashboard', name: '首页', icon: <HomeOutlined />,
        routes: [{ path: '/dashboard', name: '工作台', icon: <DashboardOutlined /> }]
      },
      ...dynamicMenus.map(m => {
        if (m.path.includes('/system')) m.icon = <SettingOutlined />;
        return m;
      }),
    ];
  }, [routers]);

  const rootPath = `/${location.pathname.split('/')[1] || 'dashboard'}`;
  
  const currentSubMenu = useMemo(() => {
    const found = fullMenuData.find(m => m.path === rootPath || (m.path === '/' && rootPath === '/dashboard'));
    if (!found) return { path: '/', routes: fullMenuData };
    return { path: found.path, routes: found.routes || [found] };
  }, [fullMenuData, rootPath]);

  const activeLevel1MenuContext = fullMenuData.find(m => m.path === rootPath || (m.path === '/' && rootPath === '/dashboard')) || fullMenuData[0];

  return (
    <div className="flex w-full h-screen overflow-hidden bg-[#f0f2f5]" style={{ display: 'flex', width: '100vw', height: '100vh', overflow: 'hidden' }}>
      {/* 1. 左侧第一级极简主导航 */}
      <div 
         className="bg-[#001529] flex flex-col items-center py-2 z-50 shrink-0 shadow-sm relative"
         style={{ width: '76px', backgroundColor: '#001529', flexShrink: 0, height: '100vh', display: 'flex', flexDirection: 'column', alignItems: 'center', zIndex: 100 }}
      >
        <div className="w-9 h-9 bg-transparent text-white rounded-lg flex items-center justify-center font-black text-xl mb-4 mt-1" style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: '24px', fontWeight: 'bold' }}>
          <span style={{ color: '#1677ff' }}>D</span><span style={{ color: '#fff' }}>G</span>
        </div>
        <div className="flex-1 w-full flex flex-col items-center gap-2 overflow-y-auto" style={{ flex: 1, width: '100%', overflowY: 'auto', display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '8px', scrollbarWidth: 'none' }}>
          {fullMenuData.map(menu => {
            const isActive = menu.path === rootPath || (menu.path === '/' && rootPath === '/dashboard');
            return (
              <div 
                key={menu.path}
                onClick={() => navigate({ to: menu.routes?.[0]?.path ? menu.routes[0].path : menu.path })}
                className={`w-[48px] h-[32px] rounded-lg flex flex-col items-center justify-center cursor-pointer transition-all duration-300 ${
                  isActive ? 'bg-[#1677ff] text-white font-medium shadow-md' : 'text-gray-400 hover:text-white hover:bg-[#ffffff14]'
                }`}
                style={{
                    width: '58px', height: '58px', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', 
                    borderRadius: '6px', cursor: 'pointer', transition: 'all 0.3s',
                    backgroundColor: isActive ? '#1677ff' : 'transparent',
                    color: isActive ? '#fff' : '#9ca3af',
                }}
              >
                <div className="text-[18px] mb-1" style={{ fontSize: '22px', marginBottom: '2px' }}>{menu.icon}</div>
                <span className="text-[11px] transform scale-90" style={{ fontSize: '13px', transform: 'scale(0.9)', whiteSpace: 'nowrap' }}>{menu.name}</span>
              </div>
            );
          })}
        </div>
      </div>

      {/* 2. 右侧 ProLayout 内容区 */}
      <div className="flex-1 w-full relative flex flex-col min-w-0 bg-[#f0f2f5]" style={{ flex: 1, height: '100vh', minWidth: 0, position: 'relative' }}>
        <ProLayout
          layout="side"
          fixedHeader={false}
          fixSiderbar={false}
          logo={null}
          title={false}
          splitMenus={false}
          siderMenuType="group"
          collapsed={collapsed}
          onCollapse={setCollapsed}
          menuHeaderRender={() => (
             <div className="h-[32px] flex items-center px-4 text-[12px] font-bold text-gray-800 tracking-wide border-b border-gray-100 bg-white shadow-sm w-full truncate" style={{ height: '38px', padding: '0 20px', fontSize: '14px' }}>
                {activeLevel1MenuContext?.name || '菜单'}
             </div>
          )}
          headerRender={false} 
          route={currentSubMenu}
          location={{ pathname: location.pathname }}
          menuItemRender={(item, dom) => (
            <a onClick={(e) => { e.preventDefault(); if (item.path) navigate({ to: item.path }); }}>{dom}</a>
          )}
          token={{
            bgLayout: '#f0f2f5',
            sider: {
              colorMenuBackground: '#fff', colorTextMenu: '#595959', colorTextMenuSelected: '#1677ff', colorBgMenuItemSelected: '#e6f4ff',
            },
            pageContainer: {
              paddingBlockPageContainerContent: 24, paddingInlinePageContainerContent: 24
            }
          }}
        >
          <div className="flex flex-col absolute inset-0 w-full h-full bg-[#f0f2f5]">
           {/* --- A. 完全自定义 Header --- */}
           <div 
             className="bg-white flex items-center justify-between px-4 shrink-0 z-10 w-full shadow-sm"
             style={{ height: '48px', backgroundColor: '#fff', display: 'flex', alignItems: 'center', justifyContent: 'space-between', padding: '0 20px', flexShrink: 0, zIndex: 10, width: '100%', boxShadow: '0 1px 2px 0 rgba(0,0,0,0.03)' }}
           >
              <div className="flex items-center gap-4" style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
                <div onClick={() => setCollapsed(!collapsed)} className="cursor-pointer" style={{ cursor: 'pointer' }}>
                    {collapsed ? <MenuUnfoldOutlined className="text-gray-500 text-lg hover:text-[#1677ff]" style={{ fontSize: '19px' }} /> : <MenuFoldOutlined className="text-gray-500 text-lg hover:text-[#1677ff]" style={{ fontSize: '19px' }} />}
                </div>
                <div className="text-gray-600 font-medium text-sm select-none" style={{ color: '#4b5563', fontSize: '16px', fontWeight: 500 }}>
                    {activeLevel1MenuContext?.name || '首页'} / {(currentSubMenu.routes?.find(r=>r.path===location.pathname) || currentSubMenu.routes?.[0])?.name || '工作台'}
                </div>
              </div>
              <div className="flex items-center gap-5 text-gray-500 text-[17px]" style={{ display: 'flex', alignItems: 'center', gap: '19px', color: '#6b7280', fontSize: '18px' }}>
                <SearchOutlined className="cursor-pointer hover:text-[#1677ff] transition-colors" style={{ cursor: 'pointer' }} />
                <LockOutlined className="cursor-pointer hover:text-[#1677ff] transition-colors" style={{ cursor: 'pointer' }} />
                <FullscreenOutlined className="cursor-pointer hover:text-[#1677ff] transition-colors" style={{ cursor: 'pointer' }} />
                <Badge dot offset={[-2, 6]}><BellOutlined className="cursor-pointer hover:text-[#1677ff] transition-colors text-[17px]" style={{ cursor: 'pointer', fontSize: '18px' }} /></Badge>
                <div className="h-4 w-[1px] bg-gray-300 mx-2" style={{ height: '14px', width: '1px', backgroundColor: '#d1d5db', margin: '0 4px' }}></div>
                <Dropdown menu={{ items: userMenuItems }} trigger={['click']}>
                  <div className="flex items-center cursor-pointer ml-1" style={{ display: 'flex', alignItems: 'center', cursor: 'pointer', marginLeft: '4px' }}>
                    <Avatar size={29} src={userInfo?.avatar || 'https://gw.alipayobjects.com/zos/antfincdn/efFD%24IOql2/weixintupian_20170331104822.jpg'} />
                    <span className="ml-2 text-sm text-gray-600 font-medium" style={{ marginLeft: '6px', fontSize: '16px', color: '#4b5563', fontWeight: 500 }}>{userInfo?.username || 'Admin'}</span>
                  </div>
                </Dropdown>
                <SettingOutlined className="cursor-pointer hover:text-[#1677ff] transition-colors" style={{ cursor: 'pointer' }} />
              </div>
            </div>

            {/* --- B. 多标签页 (Tag Views) --- */}
            <div 
              className="bg-white border-b border-gray-200 flex items-center px-4 gap-2 overflow-x-auto shrink-0 z-0 w-full"
              style={{ height: '38px', backgroundColor: '#fff', borderBottom: '1px solid #e5e7eb', display: 'flex', alignItems: 'center', padding: '0 14px', gap: '7px', overflowX: 'auto', flexShrink: 0, width: '100%' }}
            >
               <div className="bg-[#e6f4ff] text-[#1677ff] border border-[#91caff] px-4 py-1.5 rounded flex items-center text-[12px] cursor-pointer whitespace-nowrap shadow-sm" style={{ backgroundColor: '#e6f4ff', color: '#1677ff', border: '1px solid #91caff', padding: '3px 12px', borderRadius: '4px', display: 'flex', alignItems: 'center', fontSize: '14px', cursor: 'pointer', whiteSpace: 'nowrap' }}>
                 仪表盘
               </div>
            </div>
            
            {/* --- C. 内容插槽 (带水印) --- */}
            <div className="flex-1 overflow-auto p-4 w-full relative">
               <div className="absolute inset-0 pointer-events-none opacity-[0.02] z-[9999]" style={{
                  backgroundImage: "url('data:image/svg+xml,%3Csvg width=\"350\" height=\"250\" xmlns=\"http://www.w3.org/2000/svg\"%3E%3Ctext x=\"50%25\" y=\"50%25\" font-size=\"18\" fill=\"black\" font-family=\"sans-serif\" text-anchor=\"middle\" transform=\"rotate(-20 175 125)\"%3E%E8%B6%85%E7%BA%A7%E7%AE%A1%E7%90%86%E5%91%98 2026-03-12%3C/text%3E%3C/svg%3E')",
                  backgroundRepeat: 'repeat'
               }}></div>
               <div className="relative z-20 h-full">
                 <Outlet />
               </div>
            </div>
          </div>
        </ProLayout>
      </div>
    </div>
  );
};

export default BasicLayout;

