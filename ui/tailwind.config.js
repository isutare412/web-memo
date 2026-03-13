const defaultTheme = require('tailwindcss/defaultTheme')
const daisyuiThemes = require('daisyui/src/theming/themes')

const themeOverrides = { '--btn-text-case': 'none' }

/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter Variable', ...defaultTheme.fontFamily.sans],
      },
    },
  },
  plugins: [require('@tailwindcss/typography'), require('daisyui')],
  daisyui: {
    themes: [
      {
        light: { ...daisyuiThemes['light'], ...themeOverrides },
        dark: { ...daisyuiThemes['dark'], ...themeOverrides },
        cupcake: { ...daisyuiThemes['cupcake'], ...themeOverrides },
        retro: { ...daisyuiThemes['retro'], ...themeOverrides },
        halloween: { ...daisyuiThemes['halloween'], ...themeOverrides },
        sunset: { ...daisyuiThemes['sunset'], ...themeOverrides },
        dim: { ...daisyuiThemes['dim'], ...themeOverrides },
        lemonade: { ...daisyuiThemes['lemonade'], ...themeOverrides },
        bumblebee: { ...daisyuiThemes['bumblebee'], ...themeOverrides },
        corporate: { ...daisyuiThemes['corporate'], ...themeOverrides },
        forest: { ...daisyuiThemes['forest'], ...themeOverrides },
        fantasy: { ...daisyuiThemes['fantasy'], ...themeOverrides },
        business: { ...daisyuiThemes['business'], ...themeOverrides },
        coffee: { ...daisyuiThemes['coffee'], ...themeOverrides },
        nord: { ...daisyuiThemes['nord'], ...themeOverrides },
        winter: { ...daisyuiThemes['winter'], ...themeOverrides },
        valentine: { ...daisyuiThemes['valentine'], ...themeOverrides },
        autumn: { ...daisyuiThemes['autumn'], ...themeOverrides },
        garden: { ...daisyuiThemes['garden'], ...themeOverrides },
        dracula: { ...daisyuiThemes['dracula'], ...themeOverrides },
        night: { ...daisyuiThemes['night'], ...themeOverrides },
      },
    ],
  },
}
