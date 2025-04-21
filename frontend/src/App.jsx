import { useState } from "react";
import Editor from "@monaco-editor/react";
import { motion } from "framer-motion";

export default function App() {
  // Состояния для кода, вывода, уровня, аутентификации и формы
  const [code, setCode] = useState("print('Hello, world!')");
  const [output, setOutput] = useState("");
  const [level, setLevel] = useState(1);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [authMode, setAuthMode] = useState("login");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [emailError, setEmailError] = useState("");  // Состояние для ошибок почты
  const [passwordError, setPasswordError] = useState("");  // Состояние для ошибок пароля

  // Функция для обработки аутентификации
  const handleAuth = () => {
    // Проверка на пустые поля
    if (!email || !password) {
      if (!email) setEmailError("Пожалуйста, введите email.");
      if (!password) setPasswordError("Пожалуйста, введите пароль.");
      return;
    }

    // Проверка на валидность email
    const emailRegex = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$/;
    if (!emailRegex.test(email)) {
      setEmailError("Введите корректный email.");
      return;
    }

    // Если всё в порядке, аутентифицируем пользователя
    setIsAuthenticated(true);
  };

  // Функция для обработки выполнения кода
  const handleRunCode = async () => {
    try {
      const result = {
        success: code.includes("flag"),
        output: code.includes("flag")
          ? `✅ Доступ получен! Уровень ${level} пройден.`
          : "❌ Уязвимость не найдена. Попробуй ещё раз.",
      };
      setOutput(result.output);
      if (result.success) setLevel(level + 1);
    } catch (err) {
      setOutput("Ошибка выполнения кода.");
    }
  };

  // Если пользователь не аутентифицирован, показываем форму логина или регистрации
  if (!isAuthenticated) {
    return (
      <div style={{ minHeight: "100vh", display: "flex", alignItems: "center", justifyContent: "center", backgroundColor: "#f5f5f5", padding: 16 }}>
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
          style={{ width: "100%", maxWidth: 400, backgroundColor: "white", borderRadius: 16, padding: 24, boxShadow: "0 4px 12px rgba(0,0,0,0.1)" }}
        >
          <h2 style={{ fontSize: 24, fontWeight: "bold", textAlign: "center", marginBottom: 16 }}>
            {authMode === "login" ? "Вход" : "Регистрация"}
          </h2>
          {/* Ввод email с проверкой на ошибки */}
          <input
            type="email"
            placeholder="Email"
            style={{ width: "100%", border: "1px solid #ccc", borderRadius: 8, padding: 10, marginBottom: 12 }}
            value={email}
            onChange={(e) => { setEmail(e.target.value); setEmailError(""); }} // сбросить ошибку при изменении
          />
          {emailError && <p style={{ color: "red", fontSize: 12 }}>{emailError}</p>}  {/* Ошибка для email */}
          
          {/* Ввод пароля с проверкой на ошибки */}
          <input
            type="password"
            placeholder="Пароль"
            style={{ width: "100%", border: "1px solid #ccc", borderRadius: 8, padding: 10, marginBottom: 12 }}
            value={password}
            onChange={(e) => { setPassword(e.target.value); setPasswordError(""); }} // сбросить ошибку при изменении
          />
          {passwordError && <p style={{ color: "red", fontSize: 12 }}>{passwordError}</p>}  {/* Ошибка для пароля */}

          {/* Кнопка для входа или регистрации */}
          <button style={{ width: "100%", padding: 12, backgroundColor: "#007bff", color: "white", borderRadius: 8, border: "none", cursor: "pointer" }} onClick={handleAuth}>
            {authMode === "login" ? "Войти" : "Зарегистрироваться"}
          </button>

          <p style={{ fontSize: 14, textAlign: "center", marginTop: 12 }}>
            {authMode === "login" ? (
              <>
                Нет аккаунта?{' '}
                <button onClick={() => setAuthMode("register")} style={{ color: "#007bff", background: "none", border: "none", cursor: "pointer", textDecoration: "underline" }}>
                  Зарегистрироваться
                </button>
              </>
            ) : (
              <>
                Уже есть аккаунт?{' '}
                <button onClick={() => setAuthMode("login")} style={{ color: "#007bff", background: "none", border: "none", cursor: "pointer", textDecoration: "underline" }}>
                  Войти
                </button>
              </>
            )}
          </p>
        </motion.div>
      </div>
    );
  }

  // Главная часть, если пользователь аутентифицирован
  return (
    <div style={{ minHeight: "100vh", backgroundColor: "#f5f5f5", color: "#1a1a1a" }}>
      {/* Панель уровней сверху */}
      <div style={{ backgroundColor: "white", boxShadow: "0 1px 4px rgba(0,0,0,0.1)", padding: "16px 24px", display: "flex", flexWrap: "wrap", justifyContent: "flex-start" }}>
        {[...Array(20)].map((_, index) => (
          <div
            key={index}
            style={{
              width: 40,
              height: 40,
              margin: 5,
              display: "flex",
              justifyContent: "center",
              alignItems: "center",
              borderRadius: 8,
              backgroundColor: level > index ? "green" : "gray",
              color: "white",
              fontWeight: "bold",
              cursor: "pointer"
            }}
            onClick={() => setLevel(index + 1)}
          >
            {index + 1}
          </div>
        ))}
      </div>

      {/* Основной контент с задачами и редактором */}
      <div style={{ display: "flex", flexDirection: "column", gap: 24, padding: 24 }}>
        <motion.div
          initial={{ opacity: 0, x: -30 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.4 }}
          style={{ backgroundColor: "white", borderRadius: 16, padding: 16 }}
        >
          <h2 style={{ fontSize: 20, fontWeight: "bold", marginBottom: 8 }}>Уровень {level}</h2>
          <p>Найди уязвимость и получи доступ к секретным данным!</p>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, x: 30 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.4 }}
          style={{ backgroundColor: "white", borderRadius: 16, padding: 16 }}
        >
          <Editor
            height="300px"
            defaultLanguage="python"
            value={code}
            onChange={(value) => setCode(value || "")}
            theme="vs-dark"
            options={{ fontSize: 14, minimap: { enabled: false } }}
          />
          <button style={{ marginTop: 12, marginBottom: 12, padding: 12, backgroundColor: "#28a745", color: "white", border: "none", borderRadius: 8, cursor: "pointer" }} onClick={handleRunCode}>
            Отправить код
          </button>
          <textarea
            readOnly
            value={output}
            placeholder="Здесь будет результат..."
            style={{ width: "100%", height: 100, backgroundColor: "#f0f0f0", border: "1px solid #ccc", borderRadius: 8, padding: 8 }}
          />
        </motion.div>
      </div>
    </div>
  );
}
