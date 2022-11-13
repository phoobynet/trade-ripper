module.exports = {
  trailingComma: 'all',
  tabWidth: 2,
  semi: false,
  singleQuote: true,
  importOrderSeparation: true,
  importOrderSortSpecifiers: true,
  singleAttributePerLine: true,
  plugins: [require('prettier-plugin-tailwindcss')],
}