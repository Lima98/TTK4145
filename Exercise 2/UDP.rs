use std::net::UdpSocket;
use std::net::SocketAddr;
use std::net::Ipv4Addr;
use std::net::IpAddr;

fn main() -> std::io::Result<()> {

    {
    let socket = UdpSocket::bind("0.0.0.0:30000")?;
    let socket_home = SocketAddr::new(IpAddr::V4(Ipv4Addr::new(10,100,23,32)),20022);
    let mut buf = [0; 100];

    loop {
        let (num_bytes_recieved, from_who) = socket.recv_from(&mut buf)?;   
        if from_who != socket_home{
            println!("Received from server {}", num_bytes_recieved);
            println!("What dhe fuck did we receive: {:?}", &buf);

        } 
    }

    }

}