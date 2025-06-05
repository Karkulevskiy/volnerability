import { useState } from "react";
import {
  Box, Button, Container, Paper, TextField, Typography, Snackbar, Alert
} from "@mui/material";
import { motion } from "framer-motion";
import { RegisterUser } from '../api/grpcClient.js'

export default function AuthForm({ authMode, setAuthMode, onAuth }) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [emailError, setEmailError] = useState("");
  const [passwordError, setPasswordError] = useState("");
  const [authError, setAuthError] = useState("");
  const [successMessage, setSuccessMessage] = useState("");
  const [openSnackbar, setOpenSnackbar] = useState(false);

  const handleCloseSnackbar = () => {
    setOpenSnackbar(false);
  };

  async function fetchData(url, data) {
    const response = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      },
    body: JSON.stringify(data),
    });
  
    if (!response.ok) {
      const errorData = await response.text();
      throw new Error(errorData|| "Ошибка отправки http-запроса!");
    }
    const responseData = await response.json();
    return responseData
  };

  
  const handleSubmit = async () => {
    setEmailError("");
    setPasswordError("");
    setAuthError("");

    const BASE_URL = `http://localhost:8080`
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!email) setEmailError("Введите email.");
    else if (!emailRegex.test(email)) setEmailError("Некорректный email.");
    if (!password) setPasswordError("Введите пароль.");

    if (!email || !password || !emailRegex.test(email)) return;

    if (authMode === "register"){
      let reqBody = {email, password}
      try{
        await fetchData(`${BASE_URL}/register`, reqBody)
      } catch (error){
        setAuthError("Ошибка регистрации")
        console.log(error)
        return
      }
       
      setSuccessMessage("Регистрация успешно завершена! Выполняется вход...");
      setOpenSnackbar(true);
    }

    try{
      let reqBody = {email, password}
      const data = await fetchData(`${BASE_URL}/login`, reqBody) 
      localStorage.setItem("authToken", data.token)
    } catch (error) {
      setAuthError("Неверно введён логин и/или пароль")
      console.log(error)
      return
    }

    onAuth(email);
  };

  return (
    <Box display="flex" justifyContent="center" alignItems="center" minHeight="100vh">
      <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.5 }}>
        <Container maxWidth="xs">
          <Paper elevation={3} sx={{ p: 4, borderRadius: 4 }}>
            <Typography variant="h5" align="center" gutterBottom>
              {authMode === "login" ? "Вход" : "Регистрация"}
            </Typography>

            {authError && (
              <Typography color="error" align="center" sx={{ mb: 2 }}>
                {authError}
              </Typography>
            )}

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

      <Snackbar
        open={openSnackbar}
        autoHideDuration={3000}
        onClose={handleCloseSnackbar}
        anchorOrigin={{ vertical: "top", horizontal: "center" }}
      >
        <Alert onClose={handleCloseSnackbar} severity="success" sx={{ width: "100%" }}>
          {successMessage}
        </Alert>
      </Snackbar>
    </Box>
  );
}
