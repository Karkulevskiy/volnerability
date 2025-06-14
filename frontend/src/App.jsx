import { useState, useMemo, useCallback, useEffect } from "react";
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
  const [isLoading, setIsLoading] = useState(false);

  const [level, setLevel] = useState(0);
  const [levelData, setLevelData] = useState({
    id:0,
    language: "python",
    description: 'Загрузка...',
    hints: []
  })
  const [code, setCode] = useState("print('Hello, world!')");
  const [output, setOutput] = useState("");
  const [user, setUser] = useState({ email: "", completedLevels: 0, attempts: 0 });

  const theme = useMemo(() => getTheme(darkMode), [darkMode]);
  const toggleTheme = () => setDarkMode(!darkMode);

  const [currentHintIdx, setCurrentHintIdx] = useState(1);
  const [hint, setHint] = useState("");

    const fetchLevelData = useCallback(async (level) => {
    setIsLoading(true)
    try {
      if (level === 0) {
        setLevelData({
          id: 0,
          language: 'python',
          description: `Вы — хакер-новичок, нанятый таинственным заказчиком. Ваша цель — получить root-доступ к защищенной системе, которую охраняет продвинутая ИИ-система безопасности. Сначала вы должны пробраться в систему через уязвимые веб-интерфейсы, затем подорвать устойчивость серверного ПО через ошибки в управлении памятью, и, наконец, эскалировать привилегии до root-пользователя. Каждое успешно выполненное задание приближает вас к получению полного контроля над системой.`,
          hints: []
        })
        return
      }
      const token = localStorage.getItem("authToken");
      const response = await fetch(`http://localhost:8080/level?id=${level}`, {
        headers: {
          "Authorization": `Bearer ${token}`
        }
      });
      
      if (!response.ok) {
        throw new Error('Не удалось загрузить получить данные уровня');
      }
      
      const data = await response.json();
      setLevelData({
        id: data.id,
        language: data.language || 'python',
        description: data.description,
        hints: data.hints || []
      });
      setCode(data.initialCode || '');
      setCurrentHintIdx(0); // Сбрасываем индекс подсказок
      setHint("")
      setOutput("")
      
      
    } catch (error) {
      console.error('Ошибка загрузки уровня:', error);
      setLevelData(prev => ({
        ...prev,
        description: `Ошибка загрузки уровня: ${error.message}`
      }));
    } finally{
      setIsLoading(false)
    }
  }, []);

  useEffect(() => {
    fetchLevelData(level);
  }, [level, fetchLevelData]);

  const handleRunCode = async () => {
    if (level === 0){
      setLevel(1)
      return
    }
    if (code.length == 0){
      setOutput("Ошибка отправления кода(пустой ввод)")
      return
    }
    try {
      const token = localStorage.getItem("authToken");

      const res = await fetch("http://localhost:8080/submit", {
        method: "POST",
        headers: { 
          "Content-Type": "application/json",
          "Authorization": `Bearer ${token}`,
         },
        body: JSON.stringify({ levelId: level, input:code }),
      });

      if (!res.ok){
        setOutput("Ошибка отправления кода")
        throw new Error("Failed to submit")
      }
      const data = await res.json();
      setOutput(data.status || "");

      if (data.isCompleted === true) {
      setUser(prev => {
        const newCompletedLevels = Math.max(prev.completedLevels, level);
        return {
          ...prev,
          completedLevels: newCompletedLevels,
          attempts: prev.attempts + 1,
        };
      });
    } else {
      setUser(prev => ({
        ...prev,
        attempts: prev.attempts + 1
      }));
    }
    } catch (err) {
      setOutput("Ошибка при выполнении запроса.");
    }
  };

   const handleHint = useCallback(() => {
    if (levelData.hints.length == 0){
      setHint("Для этого уровня нет подсказок")
      return
    }
    setHint(levelData.hints[currentHintIdx])
    if (currentHintIdx < levelData.hints.length) {
      setCurrentHintIdx(prev => prev + 1);
    } else {
      setCurrentHintIdx(0)
    }
  });

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
          onAuth={(email, passLevels, totalAttempts) => {
            setUser({ email, completedLevels: passLevels, attempts: totalAttempts});
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
          user={user}
          handleRunCode={handleRunCode}
          handleHint={handleHint}
          hint={hint}
          darkMode={darkMode}
          toggleTheme={toggleTheme}
          setShowProfile={setShowProfile}
          levelData={levelData}
          isLoading={isLoading}
        />
      )}
    </ThemeProvider>
  );
}

