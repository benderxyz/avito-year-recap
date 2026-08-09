import { createTheme, type MantineColorsTuple, MantineProvider } from '@mantine/core';

const avito: MantineColorsTuple = [
  '#e8f7ff',
  '#cceeff',
  '#99ddff',
  '#66ccff',
  '#33bbff',
  '#00aaff',
  '#0099e6',
  '#0077b3',
  '#005580',
  '#00334d',
];

export const mantineTheme = createTheme({
  primaryColor: 'avito',
  colors: {
    avito,
  },
  fontFamily: "'Manrope', system-ui, sans-serif",
  headings: {
    fontFamily: "'Manrope', system-ui, sans-serif",
  },
  defaultRadius: 'md',
});

export { MantineProvider };
