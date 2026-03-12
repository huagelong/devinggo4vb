import React, { useState } from 'react';
import { LoginForm, ProFormText, ProFormCheckbox } from '@ant-design/pro-components';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { message } from 'antd';
import { useNavigate } from '@tanstack/react-router';
import { useAuthStore } from '../../stores/authStore';
import { login, getInfo } from '../../services/auth';
import { useTranslation } from 'react-i18next';

const LoginPage: React.FC = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const setToken = useAuthStore((state) => state.setToken);
  const setUserInfo = useAuthStore((state) => state.setUserInfo);
  const setRouters = useAuthStore((state) => state.setRouters);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (values: any) => {
    try {
      setLoading(true);
      // 调用登录接口
      const res = await login({
        username: values.username,
        password: values.password,
      });
      
      // 假设后端返回的数据包含 token
      if (res.token) {
        setToken(res.token);
        
        // 登录成功后立刻拉取个人信息与权限菜单
        const infoRes: any = await getInfo();
        if (infoRes) {
          setUserInfo(infoRes.user || infoRes);
          if (infoRes.routers) {
            setRouters(infoRes.routers);
          }
        }

        message.success('登录成功');
        navigate({ to: '/' });
      }
    } catch (error) {
      console.error(error);
      // 错误信息在 request 拦截器中已经抛出，无需在此重复
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex h-screen w-full bg-[#f0f2f5] items-center justify-center relative overflow-hidden">
      {/* Background decoration */}
      <div className="absolute top-0 right-0 -mr-20 -mt-20 w-96 h-96 rounded-full bg-blue-400 opacity-10 filter blur-3xl"></div>
      <div className="absolute bottom-0 left-0 -ml-20 -mb-20 w-96 h-96 rounded-full bg-blue-600 opacity-10 filter blur-3xl"></div>

      <div className="w-[450px] bg-white p-8 rounded-xl shadow-xl z-10 relative">
        <LoginForm
          logo="https://gw.alipayobjects.com/zos/rmsportal/KDpgvguMpGfqaHPjicRK.svg"
          title="Admin Pro UI"
          subTitle="企业级后台管理系统标准前端"
          onFinish={handleSubmit}
          submitter={{
            searchConfig: {
              submitText: t('system.login', '登 录'),
            },
            submitButtonProps: {
              loading: loading,
              size: 'large',
              style: { width: '100%', borderRadius: 6 },
            },
          }}
          containerStyle={{ paddingBottom: 0 }}
        >
          <div className="mb-8 mt-4 text-center"></div>

          <ProFormText
            name="username"
            fieldProps={{
              size: 'large',
              prefix: <UserOutlined className="text-gray-400" />,
            }}
            placeholder={'用户名 : superAdmin'}
            rules={[{ required: true, message: '请输入用户名!' }]}
          />
          <ProFormText.Password
            name="password"
            fieldProps={{
              size: 'large',
              prefix: <LockOutlined className="text-gray-400" />,
            }}
            placeholder={'密码 : admin123'}
            rules={[{ required: true, message: '请输入密码!' }]}
          />
          <div style={{ marginBlockEnd: 24 }}>
            <ProFormCheckbox noStyle name="autoLogin">
              自动登录
            </ProFormCheckbox>
            <a style={{ float: 'right' }}>忘记密码 ?</a>
          </div>
        </LoginForm>
      </div>
    </div>
  );
};

export default LoginPage;
