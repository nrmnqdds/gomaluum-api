/** @type {import('tailwindcss').Config} */
module.exports = {
	content: ["./templates/**/*.{go,js,templ,html}"],
	theme: {
		extend: {
			fontFamily: {
				poppins: ["Poppins", "sans-serif"],
			},
			colors: {
				border: "hsl(179, 30%, 18%)",
				input: "hsl(179, 30%, 18%)",
				ring: "hsl(179, 100%, 28.6%)",
				background: "hsl(240, 6%, 7%)",
				foreground: "hsl(179, 5%, 90%)",
				primary: {
					DEFAULT: "hsl(179, 100%, 28.6%)",
					foreground: "hsl(0, 0%, 100%)",
				},
				secondary: {
					DEFAULT: "hsl(179, 30%, 10%)",
					foreground: "hsl(0, 0%, 100%)",
				},
				destructive: {
					DEFAULT: "hsl(0, 100%, 30%)",
					foreground: "hsl(179, 5%, 90%)",
				},
				muted: {
					DEFAULT: "hsl(141, 30%, 15%)",
					foreground: "hsl(179, 5%, 60%)",
				},
				accent: {
					DEFAULT: "hsl(141, 30%, 15%)",
					foreground: "hsl(179, 5%, 90%)",
				},
				popover: {
					DEFAULT: "hsl(179, 50%, 5%)",
					foreground: "hsl(179, 5%, 90%)",
				},
				card: {
					DEFAULT: "hsl(179, 50%, 0%)",
					foreground: "hsl(179, 5%, 90%)",
				},
				"green-haze": {
					50: "hsl(148, 90%, 96%)",
					100: "hsl(148, 84%, 90%)",
					200: "hsl(151, 82%, 80%)",
					300: "hsl(155, 76%, 67%)",
					400: "hsl(156, 68%, 52%)",
					500: "hsl(158, 90%, 39%)",
					600: "hsl(160, 100%, 32%)",
					700: "hsl(162, 100%, 24%)",
					800: "hsl(161, 94%, 20%)",
					900: "hsl(163, 90%, 16%)",
					950: "hsl(165, 96%, 9%)",
				},
			},
			dropShadow: {
				glow: [
					"0 0px 30px rgba(255,255,255,0.8)",
					"0 0px 30px rgba(255,255,255,0.9)",
				],
			},
			borderRadius: {
				lg: "0.5rem",
				md: "calc(0.5rem - 2px)",
				sm: "calc(0.5rem - 4px)",
			},
			keyframes: {
				"accordion-down": {
					from: { height: "0" },
					to: { height: "var(--radix-accordion-content-height)" },
				},
				"accordion-up": {
					from: { height: "var(--radix-accordion-content-height)" },
					to: { height: "0" },
				},
				"infinite-scroll": {
					from: { transform: "translateX(0)" },
					to: { transform: "translateX(-100%)" },
				},
			},
			animation: {
				"accordion-down": "accordion-down 0.2s ease-out",
				"accordion-up": "accordion-up 0.2s ease-out",
				"infinite-scroll": "infinite-scroll 25s linear infinite",
			},
		},
	},
	plugins: [],
};
