import { useState } from "react";
import {
  Dialog, DialogActions, DialogContent, DialogTitle,
  TextField, Button
} from "@mui/material";

export default function ChangePasswordDialog({ open, onClose }) {
  const [oldPassword, setOldPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");

  const handleChange = () => {
    console.log("Смена пароля:", { oldPassword, newPassword });
    onClose();
  };

  return (
    <Dialog open={open} onClose={onClose}>
      <DialogTitle>Смена пароля</DialogTitle>
      <DialogContent>
        <TextField label="Старый пароль" type="password" fullWidth margin="normal"
          value={oldPassword} onChange={(e) => setOldPassword(e.target.value)} />
        <TextField label="Новый пароль" type="password" fullWidth margin="normal"
          value={newPassword} onChange={(e) => setNewPassword(e.target.value)} />
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Отмена</Button>
        <Button onClick={handleChange} variant="contained">Сменить</Button>
      </DialogActions>
    </Dialog>
  );
}

