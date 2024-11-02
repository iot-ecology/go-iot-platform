import {PageContainer} from '@ant-design/pro-components';
import {useIntl, useModel} from '@umijs/max';
import {Card, theme} from 'antd';
import React from 'react';

/**
 * 每个单独的卡片，为了复用样式抽成了组件
 * @param param0
 * @returns
 */
const InfoCard: React.FC<{
  title: string;
  index: number;
  desc: string;
  href?: string;
}> = ({title, href, index, desc}) => {
  const {useToken} = theme;
  const intl = useIntl();

  const {token} = useToken();

  return (
    <div
      style={{
        backgroundColor: token.colorBgContainer,
        boxShadow: token.boxShadow,
        borderRadius: '8px',
        fontSize: '14px',
        color: token.colorTextSecondary,
        lineHeight: '22px',
        padding: '16px 19px',
        minWidth: '220px',
        flex: 1,
      }}
    >
      <div
        style={{
          display: 'flex',
          gap: '4px',
          alignItems: 'center',
        }}
      >
        <div
          style={{
            width: 48,
            height: 48,
            lineHeight: '22px',
            backgroundSize: '100%',
            textAlign: 'center',
            padding: '8px 16px 16px 12px',
            color: '#FFF',
            fontWeight: 'bold',
            backgroundImage:
              "url('https://gw.alipayobjects.com/zos/bmw-prod/daaf8d50-8e6d-4251-905d-676a24ddfa12.svg')",
          }}
        >
          {index}
        </div>
        <div
          style={{
            fontSize: '16px',
            color: token.colorText,
            paddingBottom: 8,
          }}
        >
          {title}
        </div>
      </div>
      <div
        style={{
          fontSize: '14px',
          color: token.colorTextSecondary,
          textAlign: 'justify',
          lineHeight: '22px',
          marginBottom: 8,
        }}
      >
        {desc}
      </div>

      {href ? (
        <a href={href} target="_blank" rel="noreferrer">
          {intl.formatMessage({id: "pages.welcome.more"})} {'>'}
        </a>
      ) : (
        <div></div>
      )}
    </div>
  );
};

const Welcome: React.FC = () => {
  const {token} = theme.useToken();
  const intl = useIntl();
  const {initialState} = useModel('@@initialState');
  return (
    <PageContainer>
      <Card
        style={{
          borderRadius: 8,
        }}
        bodyStyle={{
          backgroundImage:
            initialState?.settings?.navTheme === 'realDark'
              ? 'background-image: linear-gradient(75deg, #1A1B1F 0%, #191C1F 100%)'
              : 'background-image: linear-gradient(75deg, #FBFDFF 0%, #F5F7FF 100%)',
        }}
      >
        <div
          style={{
            backgroundPosition: '100% -30%',
            backgroundRepeat: 'no-repeat',
            backgroundSize: '274px auto',
            backgroundImage:
              "url('https://gw.alipayobjects.com/mdn/rms_a9745b/afts/img/A*BuFmQqsB2iAAAAAAAAAAAAAAARQnAQ')",
          }}
        >
          <div
            style={{
              fontSize: '20px',
              color: token.colorTextHeading,
            }}
          >
            {intl.formatMessage({id: "pages.welcome.a1"})}
          </div>
          <p
            style={{
              fontSize: '14px',
              color: token.colorTextSecondary,
              lineHeight: '22px',
              marginTop: 16,
              marginBottom: 32,
              width: '65%',
            }}
          >
            {intl.formatMessage({id: "pages.welcome.a2"})}
          </p>
          <div
            style={{
              display: 'flex',
              flexWrap: 'wrap',
              gap: 16,
            }}
          >
            <InfoCard
              index={1}
              href="https://iot-dev-egi.pages.dev/"
              title={intl.formatMessage({id: "pages.welcome.a3"})}
              desc={intl.formatMessage({id: "pages.welcome.a4"})}
            />
            <InfoCard
              index={2}
              href="https://gitee.com/pychfarm_admin/go-iot-platform"
              title={intl.formatMessage({id: "pages.welcome.a5"})}
              desc={intl.formatMessage({id: "pages.welcome.a6"})}
            />

            <InfoCard
              index={3}
              href="https://gitee.com/pychfarm_admin/go-iot-platform/tree/dev/doc"
              title={intl.formatMessage({id: "pages.welcome.a7"})}
              desc={intl.formatMessage({id: "pages.welcome.a8"})}
            />
            {/*   <InfoCard
            index={4}
            title="极低的资源"
            desc="Go IoT开发平台基础运行小于1G内存（不包括数据库等服务）。"
          />*/}
          </div>
        </div>
      </Card>
    </PageContainer>
  );
};

export default Welcome;
