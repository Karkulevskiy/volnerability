import { Box, Button, Paper, Typography } from "@mui/material";

export default function ProfilePage({ user, onClose, onOpenPassword }) {
  return (
    <Box display="flex" justifyContent="center" alignItems="center" minHeight="100vh">
      <Paper elevation={3} sx={{ p: 4, borderRadius: 4, minWidth: 320 }}>
        <Typography variant="h5" gutterBottom>Профиль</Typography>
        <Typography>Email: {user.email}</Typography>
        <Typography>Пройдено уровней: {user.completedLevels}</Typography>
        <Typography>Всего попыток: {user.attempts}</Typography>

        <Button fullWidth sx={{ mt: 2 }} variant="outlined" onClick={onOpenPassword}>Сменить пароль</Button>
        <Button fullWidth sx={{ mt: 2 }} onClick={onClose}>Назад</Button>
      </Paper>
    </Box>
  );
}

