const defaultTheme = require('tailwindcss/defaultTheme');

module.exports = {
    content: [
        './ui/**/*.gohtml',
        './ui/**/*.js',
        './ui/**/*.vue'
    ],
    theme: {
        extend: {
            fontFamily: {
                sans: ['Sarabun', ...defaultTheme.fontFamily.sans]
            }
        }
    },
    plugins: [
        require('@tailwindcss/forms')
    ]
};
