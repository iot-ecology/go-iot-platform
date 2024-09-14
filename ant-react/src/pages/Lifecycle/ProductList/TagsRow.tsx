import { Space, Tag, Tooltip } from 'antd';

const MAX_TAGS_DISPLAYED = 3; // 最多显示的标签数量

const RenderTags = (tagString: string) => {
  const tags = tagString
    ?.split(',')
    .map((tag: string) => tag.trim())
    .filter((tag: string) => tag);
  const visibleTags = tags.slice(0, MAX_TAGS_DISPLAYED); // 显示的标签
  const hiddenTags = tags.slice(MAX_TAGS_DISPLAYED); // 隐藏的标签

  if (tags.length > MAX_TAGS_DISPLAYED) {
    return (
      <Space>
        {visibleTags.map((tag: string, index: number) => (
          <Tag key={index}>{tag}</Tag>
        ))}
        <Tooltip title={hiddenTags.join(', ')}>
          <Tag color="blue">查看更多</Tag>
        </Tooltip>
      </Space>
    );
  } else {
    return (
      <Space>
        {tags.map((tag: string, index: number) => (
          <Tag key={index}>{tag}</Tag>
        ))}
      </Space>
    );
  }
};

export default RenderTags;
