import {EditorView} from "@codemirror/view";

export const MY_THEME = EditorView.theme(
    {
        // 输入的字体颜色
        "&": {
            color: "#0052D9",
            backgroundColor: "#FFFFFF",
        },
        ".cm-content": {
            caretColor: "#0052D9",
        },
        // 激活背景色
        ".cm-activeLine": {
            backgroundColor: "#FAFAFA",
        },
        // 激活序列的背景色
        ".cm-activeLineGutter": {
            backgroundColor: "#FAFAFA",
        },
        // 光标的颜色
        "&.cm-focused .cm-cursor": {
            borderLeftColor: "#0052D9",
        },
        // 选中的状态
        "&.cm-focused .cm-selectionBackground, ::selection": {
            backgroundColor: "#0052D9",
            color: "#FFFFFF",
        },
        // 左侧侧边栏的颜色
        ".cm-gutters": {
            backgroundColor: "#FFFFFF",
            color: "#ddd", // 侧边栏文字颜色
            border: "none",
        },
    },
    { dark: true },
);