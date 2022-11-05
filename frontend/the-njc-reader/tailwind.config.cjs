const config = {
  content: ["./src/**/*.{html,js,svelte,ts}"],

  theme: {
    extend: {},
  },

  plugins: [require("@tailwindcss/typography"), require("daisyui")],
  daisyui: {
    themes: [{
      mytheme: {
        primary: "#990011FF",
        secondary: "#FAA21C",
        accent: "#1E5DA3",
        neutral: "rgb(83, 86, 90)",
        "base-100": "#FCF6F5FF",
      },
    },
    ]
  },
};

module.exports = config;
