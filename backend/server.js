const express = require('express');
const mongoose = require('mongoose');
const bcrypt = require('bcryptjs');
const jwt = require('jsonwebtoken');
const cors = require('cors');
const bodyParser = require('body-parser');

const app = express();
const http = require('http');
const server = http.createServer(app);
const io = require('socket.io')(server, {
  cors: { origin: '*' } // Для теста
});

io.on('connection', (socket) => {
  console.log('Client connected');
});

// В роуте добавления доната (/api/donate), после await donation.save():
io.emit('newDonation', donation); // Отправка события всем клиентам

// Замени app.listen на server.listen:
server.listen(PORT, () => console.log(`Server running on port ${PORT}`));
const PORT = process.env.PORT || 5000;
const JWT_SECRET = 'd7ae7574ce458242d8167e71db02a043e704ac3e0408f1313047d9d441f4ecb794c657fb3f33323bb0901c56d3a973e150a6ac797c211f18061a81e75aa564ee'; // Замени на свой секрет (генерируй рандомный, например, через онлайн-генератор)

// Подключи MongoDB (используй credentials из Docker)
mongoose.connect('mongodb://admin:prefectdinorah@localhost:27017/donationdb?authSource=admin', {
  useNewUrlParser: true,
  useUnifiedTopology: true,
}).then(() => console.log('MongoDB connected')).catch(err => console.log(err));

app.use(cors());
app.use(bodyParser.json());

// Модель пользователя (стример)
const UserSchema = new mongoose.Schema({
  username: { type: String, unique: true, required: true },
  password: { type: String, required: true },
});
const User = mongoose.model('User', UserSchema);

// Модель доната
const DonationSchema = new mongoose.Schema({
  streamer: { type: String, required: true },
  donorName: { type: String, required: true },
  amount: { type: Number, required: true },
  currency: { type: String, required: true },
  message: { type: String },
  timestamp: { type: Date, default: Date.now },
});
const Donation = mongoose.model('Donation', DonationSchema);

// Регистрация стримера
app.post('/api/register', async (req, res) => {
  const { username, password } = req.body;
  try {
    const hashedPw = await bcrypt.hash(password, 10);
    const user = new User({ username, password: hashedPw });
    await user.save();
    res.status(201).json({ message: 'User registered' });
  } catch (err) {
    res.status(400).json({ error: 'Username taken or error' });
  }
});

// Логин стримера
app.post('/api/login', async (req, res) => {
  const { username, password } = req.body;
  const user = await User.findOne({ username });
  if (!user || !(await bcrypt.compare(password, user.password))) {
    return res.status(401).json({ error: 'Invalid credentials' });
  }
  const token = jwt.sign({ username }, JWT_SECRET, { expiresIn: '1h' });
  res.json({ token });
});

// Middleware для проверки токена
const authMiddleware = (req, res, next) => {
  const token = req.headers.authorization?.split(' ')[1];
  if (!token) return res.status(401).json({ error: 'No token' });
  try {
    const decoded = jwt.verify(token, JWT_SECRET);
    req.user = decoded;
    next();
  } catch (err) {
    res.status(401).json({ error: 'Invalid token' });
  }
};

// Получить донаты стримера
app.get('/api/donations', authMiddleware, async (req, res) => {
  const donations = await Donation.find({ streamer: req.user.username }).sort({ timestamp: -1 });
  res.json(donations);
});

// Добавить донат
app.post('/api/donate', async (req, res) => {
  const { streamer, donorName, amount, currency, message } = req.body;
  const user = await User.findOne({ username: streamer });
  if (!user) return res.status(404).json({ error: 'Streamer not found' });
  
  const donation = new Donation({ streamer, donorName, amount, currency, message });
  await donation.save();
  res.status(201).json({ message: 'Donation added' });
});
// Временный роут удаления пользователя по имени (незащищённый)
app.delete('/api/user/:username', async (req, res) => {
  try {
    const deleted = await User.findOneAndDelete({ username: req.params.username });
    if (!deleted) return res.status(404).json({ error: 'User not found' });
    await Donation.deleteMany({ streamer: req.params.username });
    res.json({ message: 'User deleted' });
  } catch (err) {
    res.status(500).json({ error: 'Server error' });
  }
});

const { exec } = require('child_process');

// Статус сервисов (проверка, запущены ли)
app.get('/api/status', (req, res) => {
  // Проверка бэка (всегда онлайн, если запрос дошёл)
  const backendStatus = 'online';
  // Проверка фронта (пинг на порт 5173, но на сервере используй curl)
  exec('curl -s http://localhost:5173', (err) => {
    const frontendStatus = err ? 'offline' : 'online';
    res.json({ backend: backendStatus, frontend: frontendStatus });
  });
});

// Вкл/выкл (тестово, через shell; осторожно!)
app.post('/api/control', (req, res) => {
  const { service, action } = req.body; // service: 'backend' или 'frontend', action: 'start' или 'stop'
  let cmd;
  if (service === 'backend') {
    cmd = action === 'start' ? 'node server.js &' : 'pkill -f "node server.js"'; // & для фона
  } else if (service === 'frontend') {
    cmd = action === 'start' ? 'cd frontend && npm run dev -- --host &' : 'pkill -f "vite"';
  } else {
    return res.status(400).json({ error: 'Invalid service' });
  }
  exec(cmd, (err, stdout, stderr) => {
    if (err) return res.status(500).json({ error: stderr });
    res.json({ message: `${service} ${action}ed` });
  });
});

// Статус сервисов (проверка, запущены ли)
app.get('/api/status', (req, res) => {
  // Проверка бэка (всегда онлайн, если запрос дошёл)
  const backendStatus = 'online';
  // Проверка фронта (пинг на порт 5173, но на сервере используй curl)
  exec('curl -s http://localhost:5173', (err) => {
    const frontendStatus = err ? 'offline' : 'online';
    res.json({ backend: backendStatus, frontend: frontendStatus });
  });
});

// Вкл/выкл (тестово, через shell; осторожно!)
app.post('/api/control', (req, res) => {
  const { service, action } = req.body; // service: 'backend' или 'frontend', action: 'start' или 'stop'
  let cmd;
  if (service === 'backend') {
    cmd = action === 'start' ? 'node server.js &' : 'pkill -f "node server.js"'; // & для фона
  } else if (service === 'frontend') {
    cmd = action === 'start' ? 'cd frontend && npm run dev -- --host &' : 'pkill -f "vite"';
  } else {
    return res.status(400).json({ error: 'Invalid service' });
  }
  exec(cmd, (err, stdout, stderr) => {
    if (err) return res.status(500).json({ error: stderr });
    res.json({ message: `${service} ${action}ed` });
  });
});

app.listen(PORT, () => console.log(`Server running on port ${PORT}`));