using System;
using System.Net.Sockets;
using System.Text;
using System.Numerics;
using System.Security.Cryptography;
using System.IO;
using System.Buffers;
using System.Linq;
using Cryptography.EllipticCurveCryptography;

namespace Cryptography.Netwprking
{
    public class SecureSocket : IDisposable
    {
        private Socket socket;

        //System.Buffer

        MemoryPool<byte> mem_pool;

        private const int key_size = 128;
        private byte[] initialization_vector = new byte[] { 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97 };

        //public key components
        private const int pub_key_bytes = 32;

        private Curves default_curve = Curves.microsoft_160;

        public SecureSocket(Socket socket)
        {
            this.socket = socket;
            mem_pool = MemoryPool<byte>.Shared;
        }

        public void Dispose()
        {
            close();
            mem_pool.Dispose();
        }

        public void close()
        {
            socket.Close();
        }
        ~SecureSocket()
        {
            Dispose();
        }

        private SymmetricAlgorithm createEncryptionObj(string key)
        {
            SymmetricAlgorithm algorithm = Aes.Create();
            algorithm.BlockSize = key_size;
            algorithm.Key = Encoding.UTF8.GetBytes(key);
            algorithm.IV = initialization_vector;

            return algorithm;
        }

        private Span<byte> AESencrypt(Span<byte> data, string key)
        {
            SymmetricAlgorithm symmetric_algorithm = createEncryptionObj(key);

            using (MemoryStream memory_stream = new MemoryStream())
            {
                using (CryptoStream crypto_stream = new CryptoStream(memory_stream, symmetric_algorithm.CreateEncryptor(), CryptoStreamMode.Write))
                {
                    crypto_stream.Write(data);
                }
                return memory_stream.ToArray();
            }
        }

        private Span<byte> AESdecrypt(Span<byte> ciphertext, string key)
        {
            SymmetricAlgorithm symmetric_algorithm = createEncryptionObj(key);
            MemoryStream output_stream = new MemoryStream();

            using (MemoryStream memory_stream = new MemoryStream(ciphertext.ToArray()))
            {
                using (CryptoStream crypto_stream = new CryptoStream(memory_stream, symmetric_algorithm.CreateDecryptor(), CryptoStreamMode.Read))
                {
                    crypto_stream.CopyTo(output_stream);
                }
            }
            return output_stream.ToArray();
        }

        private Span<byte> generatePayload(Coordinate public_key)
        {
            Span<byte> x = public_key.x.ToByteArray();
            Span<byte> y = public_key.y.ToByteArray();

            byte[] x_padded = new byte[32];
            x.CopyTo(x_padded);

            byte[] y_padded = new byte[32];
            y.CopyTo(y_padded);

            Span<byte> payload = x_padded.Concat(y_padded).ToArray();
            return payload;
        }

        private Coordinate decodePayload(Span<byte> payload)
        {
            BigInteger x = new BigInteger(payload.Slice(0, pub_key_bytes));
            BigInteger y = new BigInteger(payload.Slice(pub_key_bytes, pub_key_bytes));

            Coordinate public_key = new Coordinate(x, y);
            return public_key;
        }

        private string sendHandshake()
        {
            KeyPair key_pair = new KeyPair(default_curve);
            ECC ecc = new ECC(default_curve);

            Span<byte> payload = generatePayload(key_pair.public_component);
            sendRaw(payload);

            Coordinate public_key = decodePayload(recvRaw(pub_key_bytes * 2));
            string key = ecc.ECDH(key_pair.private_component, public_key);

            return key;
        }

        private string recvHandshake()
        {
            KeyPair key_pair = new KeyPair(default_curve);
            ECC ecc = new ECC(default_curve);

            Coordinate public_key = decodePayload(recvRaw(pub_key_bytes * 2));
            string key = ecc.ECDH(key_pair.private_component, public_key);

            Span<byte> payload = generatePayload(key_pair.public_component);
            sendRaw(payload);

            return key;
        }


        private int sendRaw(Span<byte> data)
        {
            return socket.Send(data);
        }

        private Span<byte> recvRaw(int buffsize)
        {
            Span<byte> data_recieved;
            using (IMemoryOwner<byte> buffer = mem_pool.Rent(buffsize))
            {
                int bytes_received = socket.Receive(buffer.Memory.Span);
                data_recieved = buffer.Memory.Span.Slice(0, bytes_received);
            }

            return data_recieved;
        }

        public void sendArbitrary(Span<byte> data)
        {
            using (BinaryWriter stream = new BinaryWriter(new NetworkStream(socket)))
            {
                stream.Write(data.Length);
                stream.Write(data);
            }
        }

        public Span<byte> recvArbitrary()
        {
            using (BinaryReader stream = new BinaryReader(new NetworkStream(socket)))
            {
                int content_length = stream.ReadInt32();
                return recvRaw(content_length);
            }
        }
        public void secureSend(Span<byte> data)
        {
            string encryption_key = sendHandshake();
            Span<byte> encrypted_data = AESencrypt(data, encryption_key);

            sendArbitrary(encrypted_data);
        }

        public Span<byte> secureRecv()
        {
            string encryption_key = recvHandshake();
            Span<byte> encrypted_data = recvArbitrary();

            Span<byte> decrypted_data = AESdecrypt(encrypted_data, encryption_key);
            return decrypted_data;
        }

    }
}
