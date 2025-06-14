import { useState } from "react";
import {
  Dialog, DialogActions, DialogContent, DialogTitle,
  TextField, Button
} from "@mui/material";

export default function ChangePasswordDialog({ open, onClose }) {
  const [oldPassword, setOldPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");

  const handleChange = async () => {
    try{
      const token = localStorage.getItem("authToken");
      let reqBody = {
        "oldPassword": oldPassword,
        "newPassword": newPassword
      }

      const response = await fetch("http://localhost:8080/changePassword", {
        headers: {
          "Authorization": `Bearer ${token}`
        },
        method: "POST",
        body: JSON.stringify(reqBody)
      });

      if (!response.ok) {
        throw new Error('Не удалось сменить пароль');
      }
      
      const data = await response.json();
      if (data.status === 200){
        localStorage.setItem("authToken", data.token)
      }
    } catch(error) {
      console.log(error)
    }
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

