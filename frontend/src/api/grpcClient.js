import { createClient } from "@connectrpc/connect";
import { Auth } from "../proto/auth_pb";
import { createConnectTransport } from "@connectrpc/connect-web";

const transport = createConnectTransport({
  baseUrl: "http://localhost:8085",
});

export const  RegisterUser = async (email, password) => {
  return new Promise((resolve, reject) => {
    const client = createClient(Auth, transport)
    const res = client.login({email: email, password: paw})
    console.log(res) 
    return res
  })
}
// const client = new AuthClient(
//   'http://localhost:8085', // Адрес вашего gRPC-сервера
//   null, // credentials (обычно null)
//   { 'withCredentials': false } // Опции (отключено для локальной разработки)
// );

// export const register = (email, password) => {
//   return new Promise((resolve, reject) => {
//     const request = new proto.auth.RegisterRequest();
//     request.setEmail(email);
//     request.setPassword(password);

//     grpc.unary(client.register, {
//       request,
//       host: 'http://localhost:8085',
//       onEnd: (response) => {
//         if (response.status === grpc.Code.OK) {
//           resolve({
//             userId: response.message.getUserId()
//           });
//         } else {
//           reject(new Error(response.statusMessage || 'Registration failed'));
//         }
//       }
//     });
//   });
// };

// export const login = (email, password) => {
//   return new Promise((resolve, reject) => {
//     const request = new LoginRequest();
//     request.setEmail(email);
//     request.setPassword(password);

//     grpc.unary(client.login, {
//       request,
//       host: 'http://localhost:8085',
//       onEnd: (response) => {
//         if (response.status === grpc.Code.OK) {
//           const res = response.message;
//           resolve({
//             token: res.getToken()
//           });
//         } else {
//           reject(new Error(response.statusMessage || 'Login failed'));
//         }
//       }
//     });
//   });
// };

// // Экспортируем клиент для возможного расширения
// export default client;