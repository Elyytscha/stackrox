/*

Tailwind - The Utility-First CSS Framework

A project by Adam Wathan (@adamwathan), Jonathan Reinink (@reinink),
David Hemphill (@davidhemphill) and Steve Schoger (@steveschoger).

Welcome to the Tailwind config file. This is where you can customize
Tailwind specifically for your project. Don't be intimidated by the
length of this file. It's really just a big JavaScript object and
we've done our very best to explain each section.

View the full documentation at https://tailwindcss.com.

*/

const customFormsPlugin = require('@tailwindcss/custom-forms');
const getGradientClasses = require('./plugins/gradient');
const getGridClasses = require('./plugins/grid');
const getObjectFitClasses = require('./plugins/object-fit');
const getOrderClasses = require('./plugins/order');
const getColumnClasses = require('./plugins/columns');

const textBase = '16';

function remCalc(pixel) {
    return `${pixel / textBase}rem`;
}

module.exports = {
    important: true,
    theme: {
        screens: {
            sm: remCalc('576'),
            md: remCalc('992'),
            lg: remCalc('1250'),
            xl: remCalc('1440'),
            xxl: remCalc('1812'),
            xxxl: remCalc('2125'),
        },
        colors: {
            transparent: 'transparent',

            'base-0': 'var(--base-0)',
            'base-100': 'var(--base-100)',
            'base-200': 'var(--base-200)',
            'base-300': 'var(--base-300)',
            'base-400': 'var(--base-400)',
            'base-500': 'var(--base-500)',
            'base-600': 'var(--base-600)',
            'base-700': 'var(--base-700)',
            'base-800': 'var(--base-800)',
            'base-900': 'var(--base-900)',

            'primary-100': 'var(--primary-100)',
            'primary-200': 'var(--primary-200)',
            'primary-300': 'var(--primary-300)',
            'primary-400': 'var(--primary-400)',
            'primary-500': 'var(--primary-500)',
            'primary-600': 'var(--primary-600)',
            'primary-700': 'var(--primary-700)',
            'primary-800': 'var(--primary-800)',
            'primary-900': 'var(--primary-900)',

            'secondary-100': 'var(--secondary-100)',
            'secondary-200': 'var(--secondary-200)',
            'secondary-300': 'var(--secondary-300)',
            'secondary-400': 'var(--secondary-400)',
            'secondary-500': 'var(--secondary-500)',
            'secondary-600': 'var(--secondary-600)',
            'secondary-700': 'var(--secondary-700)',
            'secondary-800': 'var(--secondary-800)',
            'secondary-900': 'var(--secondary-900)',

            'tertiary-100': 'var(--tertiary-100)',
            'tertiary-200': 'var(--tertiary-200)',
            'tertiary-300': 'var(--tertiary-300)',
            'tertiary-400': 'var(--tertiary-400)',
            'tertiary-500': 'var(--tertiary-500)',
            'tertiary-600': 'var(--tertiary-600)',
            'tertiary-700': 'var(--tertiary-700)',
            'tertiary-800': 'var(--tertiary-800)',
            'tertiary-900': 'var(--tertiary-900)',

            'accent-100': 'var(--accent-100)',
            'accent-200': 'var(--accent-200)',
            'accent-300': 'var(--accent-300)',
            'accent-400': 'var(--accent-400)',
            'accent-500': 'var(--accent-500)',
            'accent-600': 'var(--accent-600)',
            'accent-700': 'var(--accent-700)',
            'accent-800': 'var(--accent-800)',
            'accent-900': 'var(--accent-900)',

            'success-100': 'var(--success-100)',
            'success-200': 'var(--success-200)',
            'success-300': 'var(--success-300)',
            'success-400': 'var(--success-400)',
            'success-500': 'var(--success-500)',
            'success-600': 'var(--success-600)',
            'success-700': 'var(--success-700)',
            'success-800': 'var(--success-800)',
            'success-900': 'var(--success-900)',

            'warning-100': 'var(--warning-100)',
            'warning-200': 'var(--warning-200)',
            'warning-300': 'var(--warning-300)',
            'warning-400': 'var(--warning-400)',
            'warning-500': 'var(--warning-500)',
            'warning-600': 'var(--warning-600)',
            'warning-700': 'var(--warning-700)',
            'warning-800': 'var(--warning-800)',
            'warning-900': 'var(--warning-900)',

            'caution-100': 'var(--caution-100)',
            'caution-200': 'var(--caution-200)',
            'caution-300': 'var(--caution-300)',
            'caution-400': 'var(--caution-400)',
            'caution-500': 'var(--caution-500)',
            'caution-600': 'var(--caution-600)',
            'caution-700': 'var(--caution-700)',
            'caution-800': 'var(--caution-800)',
            'caution-900': 'var(--caution-900)',

            'alert-100': 'var(--alert-100)',
            'alert-200': 'var(--alert-200)',
            'alert-300': 'var(--alert-300)',
            'alert-400': 'var(--alert-400)',
            'alert-500': 'var(--alert-500)',
            'alert-600': 'var(--alert-600)',
            'alert-700': 'var(--alert-700)',
            'alert-800': 'var(--alert-800)',
            'alert-900': 'var(--alert-900)',
        },
        letterSpacing: {
            tight: `-${remCalc('1')}`,
            normal: '0',
            wide: remCalc('.5'),
            widest: remCalc('1'),
        },
        fontFamily: {
            condensed: ['Open Sans Condensed', 'sans-serif'],
            sans: [
                'Open Sans',
                '-apple-system',
                'BlinkMacSystemFont',
                'Segoe UI',
                'Roboto',
                'Oxygen',
                'Ubuntu',
                'Cantarell',
                'Fira Sans',
                'Droid Sans',
                'Helvetica Neue',
                'sans-serif',
            ],
            serif: [
                'Constantia',
                'Lucida Bright',
                'Lucidabright',
                'Lucida Serif',
                'Lucida',
                'DejaVu Serif',
                'Bitstream Vera Serif',
                'Liberation Serif',
                'Georgia',
                'serif',
            ],
            mono: ['Menlo', 'Monaco', 'Consolas', 'Liberation Mono', 'Courier New', 'monospace'],
        },
        fontSize: {
            '2xs': remCalc('9'),
            xs: remCalc('10'),
            sm: remCalc('11'),
            base: remCalc('12'),
            lg: remCalc('13'),
            xl: remCalc('14'),
            '2xl': remCalc('16'),
            '3xl': remCalc('20'),
            '4xl': remCalc('24'),
            '5xl': remCalc('30'),
            '6xl': remCalc('40'),
        },
        fontWeight: {
            '100': 100,
            '200': 200,
            '300': 300,
            '400': 400,
            '500': 500,
            '600': 600,
            '700': 700,
            '800': 800,
            '900': 900,
        },
        lineHeight: {
            none: 1,
            tight: 1.25,
            normal: 1.5,
            loose: 2,
        },
        textColor: (theme) => theme('colors'),
        backgroundColor: (theme) => theme('colors'),
        borderWidth: {
            default: '1px',
            '0': '0',
            '2': '2px',
            '3': '3px',
            '4': '4px',
            '8': '8px',
        },
        borderColor: (theme) => {
            const colors = theme('colors');
            return {
                default: colors['grey-light'],
                ...colors,
            };
        },
        borderRadius: {
            none: '0',
            sm: remCalc('2'),
            default: remCalc('4'),
            lg: remCalc('8'),
            full: '9999px',
        },
        width: {
            auto: 'auto',
            px: '1px',
            '1': '0.25rem',
            '2': '0.5rem',
            '3': '0.75rem',
            '4': '1rem',
            '5': '1.25rem',
            '6': '1.5rem',
            '7': '1.75rem',
            '8': '2rem',
            '10': '2.5rem',
            '12': '3rem',
            '16': '4rem',
            '18': '4.25rem',
            '20': '4.5rem',
            '24': '6rem',
            '32': '8rem',
            '36': '9rem',
            '43': '10.875rem',
            '48': '12rem',
            '55': '13.75rem',
            '64': '16rem',
            '120': '46rem',
            '1/2': '50%',
            '1/3': '33.33333%',
            '2/3': '66.66667%',
            '1/4': '25%',
            '3/4': '75%',
            '1/5': '20%',
            '2/5': '40%',
            '3/5': '60%',
            '4/5': '80%',
            '1/6': '16.66666%',
            '1/7': '14.28671%',
            '1/8': '12.5%',
            '1/9': '11.11111%',
            '1/10': '10%',
            '1/12': '8.33333%',
            '5/6': '83.33333%',
            '9/10': '93%',
            full: '100%',
            screen: '100vw',
        },
        height: {
            auto: 'auto',
            px: '1px',
            '1': '0.25rem',
            '2': '0.5rem',
            '3': '0.75rem',
            '4': '1rem',
            '5': '1.25rem',
            '6': '1.5rem',
            '7': '1.75rem',
            '8': '2rem',
            '9': '2.25rem',
            '10': '2.5rem',
            '12': '3rem',
            '14': '3.5rem',
            '16': '4rem',
            '18': '4.25rem',
            '20': '4.5rem',
            '24': '6rem',
            '32': '8rem',
            '43': '10.875rem',
            '48': '12rem',
            '55': '13.75rem',
            '64': '16rem',
            '72': '20rem',
            '120': '46rem',
            full: '100%',
            screen: '100vh',
        },
        minWidth: {
            '0': '0',
            '1': '0.25rem',
            '2': '0.5rem',
            '3': '0.75rem',
            '4': '1rem',
            '5': '1.25rem',
            '6': '1.5rem',
            '7': '1.75rem',
            '8': '2rem',
            '10': '2.5rem',
            '12': '3rem',
            '16': '4rem',
            '18': '4.25rem',
            '20': '4.5rem',
            '24': '6rem',
            '32': '8rem',
            '43': '10.875rem',
            '48': '12rem',
            '55': '13.75rem',
            '64': '16rem',
            '72': '18rem',
            '108': '26.875rem',
            '120': '46rem',
            '1/2': '50%',
            '1/3': '33.33333%',
            '2/3': '66.66667%',
            '1/4': '25%',
            '3/4': '75%',
            '1/5': '20%',
            '2/5': '40%',
            '3/5': '60%',
            '4/5': '80%',
            '1/6': '16.66667%',
            '1/8': '12.5%',
            '5/6': '83.33333%',
            full: '100%',
            fit: 'fit-content',
            min: 'min-content',
            max: 'max-content',
        },
        minHeight: {
            '0': '0',
            '1': '0.25rem',
            '2': '0.5rem',
            '3': '0.75rem',
            '4': '1rem',
            '5': '1.25rem',
            '6': '1.5rem',
            '7': '1.75rem',
            '8': '2rem',
            '9': '2.25rem',
            '10': '2.5rem',
            '12': '3rem',
            '14': '3.25rem',
            '16': '4rem',
            '18': '4.25rem',
            '20': '4.5rem',
            '24': '6rem',
            '32': '8rem',
            '43': '10.875rem',
            '48': '12rem',
            '55': '13.75rem',
            '64': '16rem',
            '120': '46rem',
            '1/2': '50%',
            '1/3': '33.33333%',
            '2/3': '66.66667%',
            '1/4': '25%',
            '3/4': '75%',
            '1/5': '20%',
            '2/5': '40%',
            '3/5': '60%',
            '4/5': '80%',
            '1/6': '16.66667%',
            '1/8': '12.5%',
            '5/6': '83.33333%',
            full: '100%',
            screen: '100vh',
        },
        maxWidth: {
            '1': '0.25rem',
            '2': '0.5rem',
            '3': '0.75rem',
            '4': '1rem',
            '5': '1.25rem',
            '6': '1.5rem',
            '7': '1.75rem',
            '8': '2rem',
            '10': '2.5rem',
            '12': '3rem',
            '16': '4rem',
            '18': '4.25rem',
            '20': '4.5rem',
            '24': '6rem',
            '32': '8rem',
            '43': '10.875rem',
            '48': '12rem',
            '55': '13.75rem',
            '64': '16rem',
            '120': '46rem',
            '1/2': '50%',
            '1/3': '33.33333%',
            '2/3': '66.66667%',
            '1/4': '25%',
            '3/4': '75%',
            '1/5': '20%',
            '2/5': '40%',
            '3/5': '60%',
            '4/5': '80%',
            '1/6': '16.66667%',
            '1/8': '12.5%',
            '5/6': '83.33333%',
            xs: '20rem',
            sm: '30rem',
            md: '40rem',
            lg: '50rem',
            xl: '60rem',
            '2xl': '70rem',
            '3xl': '80rem',
            '4xl': '90rem',
            '5xl': '100rem',
            full: '100%',
            fit: 'fit-content',
            min: 'min-content',
            max: 'max-content',
        },
        maxHeight: {
            '100': '25rem',
            full: '100%',
            screen: '100vh',
        },
        padding: {
            px: '1px',
            '0': '0',
            '1': '0.25rem',
            '2': '0.5rem',
            '3': '0.75rem',
            '4': '1rem',
            '5': '1.25rem',
            '6': '1.5rem',
            '7': '1.75rem',
            '8': '2rem',
            '10': '2.5rem',
            '12': '3rem',
        },
        margin: {
            auto: 'auto',
            px: '1px',
            '0': '0',
            '1': '0.25rem',
            '2': '0.5rem',
            '3': '0.75rem',
            '4': '1rem',
            '5': '1.25rem',
            '6': '1.5rem',
            '7': '1.75rem',
            '8': '2rem',
            '-px': '-1px',
            '-0': '-0',
            '-1': '-0.25rem',
            '-2': '-0.5rem',
            '-3': '-0.75rem',
            '-4': '-1rem',
            '-6': '-1.5rem',
            '-8': '-2rem',
        },
        boxShadow: {
            default: '0 2px 8px 0 hsla(0, 0%, 0%, 0.14)',
            md: '0 8px 8px 0 hsla(0, 0%, 0%, 0.04), 0 2px 4px 0 hsla(0, 0%, 0%, 0.17)',
            lg: '0 8px 8px 0 hsla(0, 0%, 0%, 0.04), 0 2px 4px 0 hsla(0, 0%, 0%, 0.17)',
            inner: 'inset 0 0px 8px 0 hsla(0, 0%, 0%, .25)',
            none: 'none',
        },
        zIndex: {
            auto: 'auto',
            '-1': -1,
            '0': 0,
            '1': 1,
            '10': 10,
            '20': 20,
            '30': 30,
            '40': 40,
            '50': 50,
            '60': 60,
            '70': 70,
        },
        fill: {
            current: 'currentColor',
        },
        stroke: {
            current: 'currentColor',
        },
    },
    variants: {
        appearance: ['responsive', 'before', 'after'],
        backgroundAttachment: ['responsive', 'before', 'after'],
        backgroundColor: ['responsive', 'before', 'after', 'hover', 'focus'],
        backgroundPosition: ['responsive', 'before', 'after'],
        backgroundRepeat: ['responsive', 'before', 'after'],
        backgroundSize: ['responsive', 'before', 'after'],
        borderCollapse: [],
        borderColor: ['responsive', 'before', 'after', 'hover', 'focus'],
        borderRadius: ['responsive', 'before', 'after'],
        borderStyle: ['responsive', 'before', 'after'],
        borderWidth: ['responsive', 'first', 'last', 'before', 'after'],
        cursor: ['responsive', 'before', 'after'],
        display: ['responsive', 'before', 'after'],
        flexDirection: ['responsive', 'before', 'after'],
        flexWrap: ['responsive', 'before', 'after'],
        alignItems: ['responsive', 'before', 'after'],
        alignSelf: ['responsive', 'before', 'after'],
        justifyContent: ['responsive', 'before', 'after'],
        alignContent: ['responsive', 'before', 'after'],
        flex: ['responsive', 'before', 'after'],
        flexGrow: ['responsive', 'before', 'after'],
        flexShrink: ['responsive', 'before', 'after'],
        float: ['responsive', 'before', 'after'],
        fontFamily: ['responsive', 'before', 'after'],
        fontWeight: ['responsive', 'before', 'after', 'hover', 'focus'],
        height: ['responsive', 'before', 'after'],
        lineHeight: ['responsive', 'before', 'after'],
        listStylePosition: ['responsive', 'before', 'after'],
        listStyleType: ['responsive', 'before', 'after'],
        margin: ['responsive', 'first', 'last', 'before', 'after'],
        maxHeight: ['responsive', 'before', 'after'],
        maxWidth: ['responsive', 'before', 'after'],
        minHeight: ['responsive', 'before', 'after'],
        minWidth: ['responsive', 'before', 'after'],
        negativeMargin: ['responsive', 'before', 'after'],
        opacity: ['responsive', 'before', 'after'],
        outline: ['focus'],
        overflow: ['responsive', 'before', 'after'],
        padding: ['responsive', 'before', 'after'],
        pointerEvents: ['responsive', 'before', 'after'],
        position: ['responsive', 'before', 'after'],
        inset: ['responsive', 'before', 'after'],
        resize: ['responsive', 'before', 'after'],
        boxShadow: ['responsive', 'before', 'after', 'hover', 'focus'],
        fill: [],
        stroke: [],
        tableLayout: ['responsive', 'before', 'after'],
        textAlign: ['responsive', 'before', 'after'],
        textColor: ['responsive', 'before', 'after', 'hover', 'focus'],
        fontSize: ['responsive', 'before', 'after'],
        fontStyle: ['responsive', 'before', 'after', 'hover', 'focus'],
        fontSmoothing: ['responsive', 'before', 'after', 'hover', 'focus'],
        textDecoration: ['responsive', 'before', 'after', 'hover', 'focus'],
        textTransform: ['responsive', 'before', 'after', 'hover', 'focus'],
        letterSpacing: ['responsive', 'before', 'after'],
        userSelect: ['responsive', 'before', 'after'],
        verticalAlign: ['responsive', 'before', 'after'],
        visibility: ['responsive', 'before', 'after'],
        whitespace: ['responsive', 'before', 'after'],
        wordBreak: ['responsive', 'before', 'after'],
        width: ['responsive', 'before', 'after'],
        zIndex: ['responsive', 'hover', 'before', 'after'],
    },
    plugins: [
        customFormsPlugin,
        function addvariant({ addVariant, e }) {
            addVariant('before', ({ modifySelectors, separator }) => {
                modifySelectors(
                    ({ className }) => `.${e(`before${separator}${className}`)}:before`
                );
            });
            addVariant('after', ({ modifySelectors, separator }) => {
                modifySelectors(({ className }) => `.${e(`after${separator}${className}`)}:after`);
            });
        },
        getGradientClasses({
            variants: ['responsive', 'before', 'after', 'hover'],
        }),
        getObjectFitClasses({
            variants: ['responsive', 'before', 'after'],
        }),
        getGridClasses({
            gaps: {
                '0': '0',
                '1px': '1px',
                '1': '0.25rem',
                '2': '0.5rem',
                '3': '0.75rem',
                '4': '1rem',
                '5': '1.25rem',
                '6': '1.5rem',
                '7': '1.75rem',
                '8': '2rem',
                '10': '2.5rem',
                '12': '3rem',
                '16': '4rem',
            },
            variants: ['responsive', 'before', 'after'],
        }),
        getOrderClasses({
            positive: {
                '0': 0,
                '1': 1,
                '2': 2,
                '3': 3,
                '4': 4,
                '5': 5,
            },
            negative: {
                '-1': -1,
                '-2': -2,
                '-3': -3,
                '-4': -4,
                '-5': -5,
            },
            variants: ['responsive', 'before', 'after'],
        }),
        getColumnClasses({
            index: {
                '0': 0,
                '1': 1,
                '2': 2,
                '3': 3,
                '4': 4,
                '5': 5,
            },
            variants: ['responsive', 'before', 'after'],
        }),
    ],
    corePlugins: {
        container: false,
    },
};
