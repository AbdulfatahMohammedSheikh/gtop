/** @type {import('tailwindcss').Config} */
module.exports = {

    content: ["./static/**.{js,css}", "./template/**/*.html"],
    theme: {
        extend: {},
    },
    plugins: [],
}

// for using templ files
// module.exports = {
//     content: [ "./**/*.html", "./**/*.templ", "./**/*.go", ],
//     theme: { extend: {}, },
//     plugins: [],
// }
