import { GithubOutlined } from '@ant-design/icons';
import { DefaultFooter } from '@ant-design/pro-components';
import React from 'react';

const Footer: React.FC = () => {
  return (
    <DefaultFooter
      style={{
        background: 'none',
      }}
      links={[
        {
          key: 'Go IoT 开发平台',
          title: 'Go IoT开发平台',
          href: 'https://iot-dev-egi.pages.dev/',
          blankTarget: true,
        },
        {
          key: 'gitee',
          title: <GithubOutlined />,
          href: 'https://gitee.com/pychfarm_admin/go-iot-platform',
          blankTarget: true,
        },

      ]}
    />
  );
};

export default Footer;
