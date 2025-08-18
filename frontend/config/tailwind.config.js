/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        orange: {
          500: '#ff6b35',
          600: '#e55a2b',
          700: '#cc4f24',
        }
      }
    },
  },
  plugins: [],
}