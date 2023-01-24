const config = {
  content: ["./src/**/*.{html,js,svelte,ts}"],

  theme: {
    extend: {},
  },

  plugins: [require("@tailwindcss/typography"), require("daisyui")],
  daisyui: {
    themes: [{
      mytheme: {
        primary: "#990011",
        secondary: "#FAA21C",
        accent: "#437D00",
        neutral: "#53565A",
        "base-100": "#FCF6F5FF",
      },
    },
    ]
  },
};

module.exports = config;
