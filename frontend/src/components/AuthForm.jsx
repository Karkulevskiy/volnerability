import { useState } from "react";
import {
  Box, Button, Container, Paper, TextField, Typography
} from "@mui/material";
import { motion } from "framer-motion";
import { RegisterUser } from '../api/grpcClient.js'

export default function AuthForm({ authMode, setAuthMode, onAuth }) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [emailError, setEmailError] = useState("");
  const [passwordError, setPasswordError] = useState("");

  const handleSubmit =  async () => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!email) setEmailError("Введите email.");
    else if (!emailRegex.test(email)) setEmailError("Некорректный email.");
    if (!password) setPasswordError("Введите пароль.");

    if (email && password && emailRegex.test(email)) try {
      if (authMode === "login"){
        // const response = await authClient.login(new LoginRequest({email, password}))
        // if (response.token) {
        //   onAuth(email); // Успешная авторизация
        // } else {
        //   setEmailError(response.error || 'Ошибка авторизации');
        // }
      } else {
        const response = await RegisterUser(email, password)
      }

    } catch (err) {
      console.log(err.toString());
    }
    setEmailError("");
    setPasswordError("");

  };

  return (
    <Box display="flex" justifyContent="center" alignItems="center" minHeight="100vh">
      <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.5 }}>
        <Container maxWidth="xs">
          <Paper elevation={3} sx={{ p: 4, borderRadius: 4 }}>
            <Typography variant="h5" align="center" gutterBottom>
              {authMode === "login" ? "Вход" : "Регистрация"}
            </Typography>

            <TextField fullWidth label="Email" type="email" margin="normal" value={email}
              onChange={(e) => setEmail(e.target.value)} error={!!emailError} helperText={emailError} />
            <TextField fullWidth label="Пароль" type="password" margin="normal" value={password}
              onChange={(e) => setPassword(e.target.value)} error={!!passwordError} helperText={passwordError} />

            <Button fullWidth variant="contained" sx={{ mt: 2 }} onClick={handleSubmit}>
              {authMode === "login" ? "Войти" : "Зарегистрироваться"}
            </Button>

            <Button fullWidth variant="outlined" sx={{ mt: 2 }} onClick={() => onAuth("google@example.com")}>
              Войти через Google
            </Button>

            <Box textAlign="center" mt={2}>
              <Button variant="text" onClick={() => setAuthMode(authMode === "login" ? "register" : "login")}>
                {authMode === "login" ? "Нет аккаунта? Зарегистрироваться" : "Уже есть аккаунт? Войти"}
              </Button>
            </Box>
          </Paper>
        </Container>
      </motion.div>
    </Box>
  );
}

