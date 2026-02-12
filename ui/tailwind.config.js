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
      },
    ],
  },
}
