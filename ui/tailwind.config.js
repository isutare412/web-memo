const defaultTheme = require('tailwindcss/defaultTheme')

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
  plugins: [require('daisyui')],
  daisyui: {
    themes: [
      {
        bumblebee: {
          ...require('daisyui/src/theming/themes')['[data-theme=bumblebee]'],
          '--btn-text-case': 'none',
        },
        dracula: {
          ...require('daisyui/src/theming/themes')['[data-theme=dracula]'],
          '--btn-text-case': 'none',
        },
      },
    ],
  },
}
