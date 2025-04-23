// App.js
import { useState, useMemo } from "react";
import { ThemeProvider, createTheme } from "@mui/material/styles";
import {
  CssBaseline,
  IconButton,
  Container,
  Box,
  Typography,
  Button,
  TextField,
  Paper,
  Grid,
  Tabs,
  Tab,
  Divider
} from "@mui/material";
import { Brightness4, Brightness7 } from "@mui/icons-material";
import { motion } from "framer-motion";
import Editor from "@monaco-editor/react";

function AuthForm({ authMode, setAuthMode, onAuth }) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [emailError, setEmailError] = useState("");
  const [passwordError, setPasswordError] = useState("");

  const handleSubmit = () => {
    setEmailError("");
    setPasswordError("");

    if (!email) setEmailError("Введите email.");
    if (!password) setPasswordError("Введите пароль.");

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (email && !emailRegex.test(email)) setEmailError("Некорректный email.");

    if (email && password && emailRegex.test(email)) onAuth();
  };

  return (
    <Box display="flex" justifyContent="center" alignItems="center" minHeight="100vh">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
      >
        <Container maxWidth="xs">
          <Paper elevation={3} sx={{ p: 4, borderRadius: 4 }}>
            <Typography variant="h5" align="center" gutterBottom>
              {authMode === "login" ? "Вход" : "Регистрация"}
            </Typography>

            <TextField
              fullWidth
              label="Email"
              type="email"
              margin="normal"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              error={Boolean(emailError)}
              helperText={emailError}
            />

            <TextField
              fullWidth
              label="Пароль"
              type="password"
              margin="normal"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              error={Boolean(passwordError)}
              helperText={passwordError}
            />

            <Button
              fullWidth
              variant="contained"
              color="primary"
              sx={{ mt: 2 }}
              onClick={handleSubmit}
            >
              {authMode === "login" ? "Войти" : "Зарегистрироваться"}
            </Button>

            <Box textAlign="center" mt={2}>
              <Button
                variant="text"
                onClick={() => setAuthMode(authMode === "login" ? "register" : "login")}
              >
                {authMode === "login" ? "Нет аккаунта? Зарегистрироваться" : "Уже есть аккаунт? Войти"}
              </Button>
            </Box>
          </Paper>
        </Container>
      </motion.div>
    </Box>
  );
}

function MainScreen({ level, setLevel, code, setCode, output, handleRunCode }) {
  return (
  <Container maxWidth="xl" disableGutters sx={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
  {/* Верх: кнопки уровней */}
  <Box sx={{ p: 2 }}>
    <Paper elevation={3} sx={{ p: 2, borderRadius: 2 }}>
      <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
        {[...Array(20)].map((_, index) => (
          <Button
            key={index}
            size="small"
            variant={level === index + 1 ? "contained" : "outlined"}
            color={level > index ? "success" : "inherit"}
            onClick={() => setLevel(index + 1)}
          >
            {index + 1}
          </Button>
        ))}
      </Box>
    </Paper>
  </Box>

  {/* Нижняя часть */}
  <Box sx={{ flex: 1, display: 'flex', overflow: 'hidden', px: 2, pb: 2 }}>
    {/* Левая колонка: задача */}
    <Box sx={{ width: '30%', minWidth: 280, overflowY: 'auto' }}>
      <Paper
        elevation={3}
        sx={{
          p: 2,
          borderRadius: 3,
          height: '100%',
          border: '1px solid #ccc',
          backgroundColor: '#f9f9f9'
        }}
      >
        <Typography variant="h6" gutterBottom>
          Задача уровня {level}
        </Typography>
        <Typography variant="body1">
          Найди уязвимость и получи доступ к секретным данным!
        </Typography>
      </Paper>
    </Box>

    {/* Правая колонка */}
    <Box sx={{ flex: 1, ml: 2, display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
      {/* Редактор */}
      <Paper
        elevation={3}
        sx={{
          p: 2,
          borderRadius: 3,
          flex: 1,
          display: 'flex',
          flexDirection: 'column',
          overflow: 'hidden',
          backgroundColor: '#1e1e1e',
        }}
      >
        <Box sx={{ flex: 1, minHeight: 0 }}>
          <Editor
            height="100%"
            defaultLanguage="python"
            value={code}
            onChange={(val) => setCode(val || "")}
            theme="vs-dark"
            options={{ fontSize: 14, minimap: { enabled: false } }}
          />
        </Box>
        <Button
          variant="contained"
          color="success"
          sx={{ mt: 2, alignSelf: 'flex-start', borderRadius: 2 }}
          onClick={handleRunCode}
        >
          Отправить код
        </Button>
      </Paper>

      {/* Терминал */}
      <Paper
        elevation={3}
        sx={{
          p: 2,
          borderRadius: 3,
          mt: 2,
          border: '1px solid #ccc',
          backgroundColor: '#f1f1f1'
        }}
      >
        <Typography variant="subtitle1" gutterBottom>
          Результат выполнения:
        </Typography>
        <TextField
          fullWidth
          multiline
          value={output}
          InputProps={{
            readOnly: true,
            sx: { fontFamily: 'monospace' }
          }}
          rows={4}
          placeholder="Здесь будет результат..."
        />
      </Paper>
    </Box>
  </Box>
</Container>

  );
}

export default function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [authMode, setAuthMode] = useState("login");
  const [darkMode, setDarkMode] = useState(false);

  const [level, setLevel] = useState(1);
  const [code, setCode] = useState("print('Hello, world!')");
  const [output, setOutput] = useState("");

  const theme = useMemo(() => createTheme({
    palette: {
      mode: darkMode ? "dark" : "light",
    },
    typography: {
      fontFamily: "'Figtree', 'Roboto', sans-serif",
    },
  }), [darkMode]);

  const toggleTheme = () => setDarkMode(!darkMode);

  const handleRunCode = () => {
    const result = {
      success: code.includes("flag"),
      output: code.includes("flag")
        ? `✅ Доступ получен! Уровень ${level} пройден.`
        : "❌ Уязвимость не найдена. Попробуй ещё раз.",
    };
    setOutput(result.output);
    if (result.success) setLevel(level + 1);
  };

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />

      <Box sx={{ position: "fixed", top: 16, right: 16, zIndex: 1200 }}>
        <IconButton onClick={toggleTheme} color="inherit">
          {darkMode ? <Brightness7 /> : <Brightness4 />}
        </IconButton>
      </Box>

      {!isAuthenticated ? (
        <AuthForm
          authMode={authMode}
          setAuthMode={setAuthMode}
          onAuth={() => setIsAuthenticated(true)}
        />
      ) : (
        <MainScreen
          level={level}
          setLevel={setLevel}
          code={code}
          setCode={setCode}
          output={output}
          handleRunCode={handleRunCode}
        />
      )}
    </ThemeProvider>
  );
}
