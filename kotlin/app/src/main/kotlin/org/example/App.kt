package appstack.kotlin


import io.ktor.server.application.*
import io.ktor.server.engine.*
import io.ktor.server.netty.*
import io.ktor.server.engine.*
import io.ktor.server.netty.*
import io.ktor.server.application.*
import io.ktor.server.plugins.callloging.*
import io.ktor.server.routing.*
import io.ktor.server.response.*
import io.ktor.http.*
import io.ktor.server.plugins.contentnegotiation.*
import org.jetbrains.exposed.sql.*
import org.jetbrains.exposed.sql.transactions.*
import org.jetbrains.exposed.dao.*
import org.jetbrains.exposed.dao.id.*
import org.slf4j.event.*
import io.ktor.serialization.kotlinx.json.*
import kotlinx.serialization.json.Json
import kotlinx.serialization.encodeToString
import kotlinx.serialization.Serializable

@Serializable
data class ItemDto(
    val id: Int,
    val name: String,
    val description: String
)

object ItemsTable : IntIdTable("items") {
    val name = text("name")
    val description = text("description")    
}

class ItemsDao(id: EntityID<Int>) : IntEntity(id) {
    companion object : IntEntityClass<ItemsDao>(ItemsTable)
    var name by ItemsTable.name
    var description by ItemsTable.description
}

fun daoToDto(dao: ItemsDao) = ItemDto(
    dao.id.value,
    dao.name,
    dao.description
)

fun main() {
    Database.connect("jdbc:sqlite:../../items.db", driver = "org.sqlite.JDBC")

    transaction {
        SchemaUtils.create(ItemsTable)
    }

    // Ktor Server
    embeddedServer(Netty, port = 8080) {
        install(ContentNegotiation) {
            json()
        }
        
        install(CallLogging) {
            //logger = LoggerFactory.getLogger(Application::class.java)
            level = Level.DEBUG  // Adjust the level as needed
            // Optionally, filter requests to log
            //filter { call -> call.request.path().startsWith("/") }
        }

        routing {
            get("/items") {
                val items = transaction { ItemsDao.all().map(::daoToDto) }
                call.respond(items)
            }
            
            // Add more routes for POST, PUT, DELETE operations
        }
    }.start(wait = true)
}
