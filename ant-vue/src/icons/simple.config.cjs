module.exports = {
  plugins: [
    "preset-default",
    "prefixIds",
    {
      name: "removeAttrs",
      params: {
        attrs: ["stroke", "fill", "fill-rule"],
      },
    },
  ],
};
