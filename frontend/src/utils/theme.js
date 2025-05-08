import { createTheme } from "@mui/material";

export const getTheme = (darkMode) =>
  createTheme({
    palette: {
      mode: darkMode ? "dark" : "light",
    },
    typography: {
      fontFamily: "monospace, sans-serif",
      fontSize: 16,
    },
  });

