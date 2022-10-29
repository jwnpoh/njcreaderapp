const config = {
  content: ["./src/**/*.{html,js,svelte,ts}"],

  theme: {
    extend: {},
  },

  plugins: [require("@tailwindcss/typography"), require("daisyui")],
  daisyui: {
    themes: [{
      mytheme: {
        primary: "#e53935",
        secondary: "rgb(254, 219, 0)",
        accent: "#37cdbe",
        neutral: "rgb(83, 86, 90)",
        "base-100": "#efefef",
      },
    },
    ]
  },
};

module.exports = config;
