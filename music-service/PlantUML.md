@startuml
    class User {
        +ID: int
        +Email: string
        +Password: string
        +Register()
        +Login()
    }
    
    class Artist {
        +ID: int
        +Name: string
    }
    
    class Album {
        +ID: int
        +Title: string
        +ArtistID: int
    }
    
    class Track {
        +ID: int
        +Title: string
        +Duration: int
        +AlbumID: int
    }
    
    class Playlist {
        +ID : int
        +Name : string
        +UserID : int
        +AddTrack(trackID)
        +RemoveTrack(trackID)
    }
    
    class MusicService {
        +PlayTrack(trackID)
        +SearchTrack(query)
    }
    
    class PlaylistService {
        +CreatePlaylist(userID)
        +AddTrackToPlaylist(playlistID, trackID)
    }
    
    class UserService {
        +Register(email, password)
        +Login(email, password)
    }

    User "1" --> "*" Playlist
    Artist "1" --> "*" Album "1" --> "*" Track
    Playlist "*" --> "*" Track: contains
    User "*" --> "*" Track: favorites
    
    MusicService --> Track
    PlaylistService --> Playlist
    UserService --> User
@enduml
