import { useState, useMemo } from "react";
import { CssBaseline, IconButton, ThemeProvider, Box, Button } from "@mui/material";
import { Brightness4, Brightness7 } from "@mui/icons-material";
import { getTheme } from "./utils/theme";
import AuthForm from "./components/AuthForm";
import MainScreen from "./components/MainScreen";
import ProfilePage from "./components/ProfilePage";
import ChangePasswordDialog from "./components/ChangePasswordDialog";

export default function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [authMode, setAuthMode] = useState("login");
  const [darkMode, setDarkMode] = useState(false);
  const [showProfile, setShowProfile] = useState(false);
  const [showPasswordDialog, setShowPasswordDialog] = useState(false);

  const [level, setLevel] = useState(1);
  const [attempts, setAttempts] = useState(1);
  const [code, setCode] = useState("print('Hello, world!')");
  const [output, setOutput] = useState("");
  const [user, setUser] = useState({ email: "", completedLevels: 0, attempts: 0 });

  const theme = useMemo(() => getTheme(darkMode), [darkMode]);
  const toggleTheme = () => setDarkMode(!darkMode);


  const [hint, setHint] = useState("");

  const handleRunCode = async () => {
    try {
      const token = localStorage.getItem("authToken");

      const res = await fetch("/api/submit", {
        method: "POST",
        headers: { 
          "Content-Type": "application/json",
          "Authorization": `Bearer ${token}`,
         },
        body: JSON.stringify({ levelId: level, code }),
      });
      const data = await res.json();
      setOutput(data.output || "");
      setUser(prev => {
        const newAttempts = prev.attempts + 1;
        let newCompletedLevels = prev.completedLevels;
        
        if (data.success) {
          newCompletedLevels = newCompletedLevels + 1;
          setLevel((prev) => prev + 1);
        }
        
        return {
          ...prev,
          attempts: newAttempts,
          completedLevels: newCompletedLevels,
        };
      });
    } catch (err) {
      setOutput("Ошибка при выполнении запроса.");
    }
  };

  const handleHint = async () => {
    try {
      const token = localStorage.getItem("authToken");

      const res = await fetch(`0.0.0.0:8080/hint?hintId=${hint}`, {
        headers: {
          "Authorization": `Bearer ${token}`,
        }
      });
      const data = await res.json();
      setHint(data.hint || "Подсказка недоступна.");
    } catch {
      setHint("Ошибка при получении подсказки.");
    }
  };

const handleLevelChange = (newLevel) => {
    setLevel(newLevel);
    setUser(prev => {
      if (newLevel > prev.completedLevels) {
        return {
          ...prev,
          completedLevels: newLevel
        };
      }
      return prev;
    });
  };

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />

      <ChangePasswordDialog open={showPasswordDialog} onClose={() => setShowPasswordDialog(false)} />

      {showProfile ? (
        <ProfilePage
          user={user}
          onClose={() => setShowProfile(false)}
          onOpenPassword={() => setShowPasswordDialog(true)}
        />
      ) : !isAuthenticated ? (
        <AuthForm
          authMode={authMode}
          setAuthMode={setAuthMode}
          onAuth={(email) => {
            setUser({ email, completedLevels: 0, attempts: 0 });
            setIsAuthenticated(true);
          }}
        />
      ) : (
        <MainScreen
          level={level}
          setLevel={setLevel}
          code={code}
          setCode={setCode}
          output={output}
          handleRunCode={handleRunCode}
          handleHint={handleHint}
          hint={hint}
          darkMode={darkMode}
          toggleTheme={toggleTheme}
          setShowProfile={setShowProfile}
        />
      )}
    </ThemeProvider>
  );
}

