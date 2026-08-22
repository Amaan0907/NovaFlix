import React from 'react'

const Movie = ({ movie }) => {
    return (
        <div className='mb-4 w-full md:w-1/3'>
            <div className='card h-full shadow-2xs'>
                <div style={{ position: "relative" }}>
                    <img src=
                        {movie.poster_path} alt={movie.title} className='card-img-top' style={
                            {
                                objectFit:"contain",
                                height:"250px",
                                width:"100%"
                            }
                        }/>
                </div>
                <div className='card-body flex flex-col'>
                    <h5 className='card-title'>{movie.title}</h5>
                    <p className='card-text mb-2'>{movie.imdb_id}</p>

                </div>
                {movie.ranking?.ranking_name && (
                    <span className='badge bg-gray-900 m-3 p-2' style={{ fontSize:"12px" }}>
                        {movie.ranking.ranking_name}
                    </span>
                )}
            </div>
        </div>
    )
}

export default Movie
